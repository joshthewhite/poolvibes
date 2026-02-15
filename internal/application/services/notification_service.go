package services

import (
	"context"
	"fmt"
	"log/slog"
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
	slog.Info("Notification scheduler started", "interval", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run immediately on start
	s.checkAndNotify(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Notification scheduler stopped")
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
		slog.Error("Notification check error", "error", err)
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
			slog.Error("Notification: could not find user", "userID", userTasks[0].UserID, "error", err)
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

	// Email notification — claim first, send only if we won the claim
	if s.emailNotifier != nil && user.NotifyEmail && user.Email != "" {
		notif := entities.NewTaskNotification(task.ID, user.ID, "email", dueDate)
		claimed, err := s.notifRepo.Claim(ctx, notif)
		if err != nil {
			slog.Error("Email claim error", "taskID", task.ID, "error", err)
		} else if claimed {
			if err := s.emailNotifier.Send(ctx, user.Email, subject, body); err != nil {
				slog.Error("Email send error", "taskID", task.ID, "email", user.Email, "error", err)
				// Release claim so another tick can retry
				if delErr := s.notifRepo.Delete(ctx, notif.ID); delErr != nil {
					slog.Error("Error releasing email claim", "error", delErr)
				}
			} else {
				slog.Info("Email notification sent", "task", task.Name, "email", user.Email)
			}
		}
	}

	// SMS notification — claim first, send only if we won the claim
	if s.smsNotifier != nil && user.NotifySMS && user.Phone != "" {
		notif := entities.NewTaskNotification(task.ID, user.ID, "sms", dueDate)
		claimed, err := s.notifRepo.Claim(ctx, notif)
		if err != nil {
			slog.Error("SMS claim error", "taskID", task.ID, "error", err)
		} else if claimed {
			if err := s.smsNotifier.Send(ctx, user.Phone, subject, body); err != nil {
				slog.Error("SMS send error", "taskID", task.ID, "phone", user.Phone, "error", err)
				// Release claim so another tick can retry
				if delErr := s.notifRepo.Delete(ctx, notif.ID); delErr != nil {
					slog.Error("Error releasing SMS claim", "error", delErr)
				}
			} else {
				slog.Info("SMS notification sent", "task", task.Name, "phone", user.Phone)
			}
		}
	}
}
