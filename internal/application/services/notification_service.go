package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type NotificationService struct {
	taskRepo      repositories.TaskRepository
	userRepo      repositories.UserRepository
	notifRepo     repositories.TaskNotificationRepository
	emailNotifier Notifier
	smsNotifier   Notifier
	interval      time.Duration
}

func NewNotificationService(
	taskRepo repositories.TaskRepository,
	userRepo repositories.UserRepository,
	notifRepo repositories.TaskNotificationRepository,
	emailNotifier Notifier,
	smsNotifier Notifier,
	interval time.Duration,
) *NotificationService {
	return &NotificationService{
		taskRepo:      taskRepo,
		userRepo:      userRepo,
		notifRepo:     notifRepo,
		emailNotifier: emailNotifier,
		smsNotifier:   smsNotifier,
		interval:      interval,
	}
}

func (s *NotificationService) Start(ctx context.Context) {
	log.Printf("Notification scheduler started (interval: %s)", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run immediately on start
	s.checkAndNotify(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Notification scheduler stopped")
			return
		case <-ticker.C:
			s.checkAndNotify(ctx)
		}
	}
}

func (s *NotificationService) checkAndNotify(ctx context.Context) {
	today := time.Now()
	tasks, err := s.taskRepo.FindDueOnDate(ctx, today)
	if err != nil {
		log.Printf("Notification check error: %v", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	// Group tasks by user
	byUser := make(map[string][]entities.Task)
	for _, t := range tasks {
		uid := t.UserID.String()
		byUser[uid] = append(byUser[uid], t)
	}

	for _, userTasks := range byUser {
		if len(userTasks) == 0 {
			continue
		}
		user, err := s.userRepo.FindByID(ctx, userTasks[0].UserID)
		if err != nil || user == nil {
			log.Printf("Notification: could not find user %s: %v", userTasks[0].UserID, err)
			continue
		}

		for _, task := range userTasks {
			s.notifyForTask(ctx, user, &task, today)
		}
	}
}

func (s *NotificationService) notifyForTask(ctx context.Context, user *entities.User, task *entities.Task, today time.Time) {
	dueDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	subject := fmt.Sprintf("PoolVibes: Task due today — %s", task.Name)
	body := fmt.Sprintf("Your pool maintenance task \"%s\" is due today (%s).", task.Name, task.DueDate.Format("Jan 2, 2006"))
	if task.Description != "" {
		body += fmt.Sprintf("\n\nDetails: %s", task.Description)
	}

	// Email notification
	if s.emailNotifier != nil && user.NotifyEmail && user.Email != "" {
		exists, err := s.notifRepo.ExistsByTaskAndType(ctx, task.ID, "email", dueDate)
		if err != nil {
			log.Printf("Notification check error for task %s: %v", task.ID, err)
		} else if !exists {
			if err := s.emailNotifier.Send(ctx, user.Email, subject, body); err != nil {
				log.Printf("Email notification error for task %s to %s: %v", task.ID, user.Email, err)
			} else {
				notif := entities.NewTaskNotification(task.ID, user.ID, "email", dueDate)
				if err := s.notifRepo.Create(ctx, notif); err != nil {
					log.Printf("Error recording email notification: %v", err)
				}
				log.Printf("Email notification sent for task %s to %s", task.Name, user.Email)
			}
		}
	}

	// SMS notification
	if s.smsNotifier != nil && user.NotifySMS && user.Phone != "" {
		exists, err := s.notifRepo.ExistsByTaskAndType(ctx, task.ID, "sms", dueDate)
		if err != nil {
			log.Printf("Notification check error for task %s: %v", task.ID, err)
		} else if !exists {
			if err := s.smsNotifier.Send(ctx, user.Phone, subject, body); err != nil {
				log.Printf("SMS notification error for task %s to %s: %v", task.ID, user.Phone, err)
			} else {
				notif := entities.NewTaskNotification(task.ID, user.ID, "sms", dueDate)
				if err := s.notifRepo.Create(ctx, notif); err != nil {
					log.Printf("Error recording SMS notification: %v", err)
				}
				log.Printf("SMS notification sent for task %s to %s", task.Name, user.Phone)
			}
		}
	}
}
