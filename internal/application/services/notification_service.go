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
	pushNotifier  PushNotifier
	interval      time.Duration
}

func NewNotificationService(
	taskRepo repositories.TaskRepository,
	userRepo repositories.UserRepository,
	notifRepo repositories.TaskNotificationRepository,
	emailNotifier Notifier,
	smsNotifier Notifier,
	pushNotifier PushNotifier,
	interval time.Duration,
) *NotificationService {
	return &NotificationService{
		taskRepo:      taskRepo,
		userRepo:      userRepo,
		notifRepo:     notifRepo,
		emailNotifier: emailNotifier,
		smsNotifier:   smsNotifier,
		pushNotifier:  pushNotifier,
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

		s.notifyBatch(ctx, user, userTasks, today)
	}
}

// notifyBatch sends at most one notification per user per channel per day,
// batching all due tasks into a single message.
func (s *NotificationService) notifyBatch(ctx context.Context, user *entities.User, tasks []entities.Task, today time.Time) {
	dueDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	subject := fmt.Sprintf("PoolVibes: %d task(s) due today", len(tasks))
	body := formatBatchBody(tasks)

	// Email notification — claim once per user per day
	if s.emailNotifier != nil && user.NotifyEmail && user.Email != "" {
		notif := entities.NewBatchNotification(user.ID, "email", dueDate)
		claimed, err := s.notifRepo.Claim(ctx, notif)
		if err != nil {
			slog.Error("Email claim error", "userID", user.ID, "error", err)
		} else if claimed {
			if err := s.emailNotifier.Send(ctx, user.Email, subject, body); err != nil {
				slog.Error("Email send error", "userID", user.ID, "email", user.Email, "error", err)
				if delErr := s.notifRepo.Delete(ctx, notif.ID); delErr != nil {
					slog.Error("Error releasing email claim", "error", delErr)
				}
			} else {
				slog.Info("Email notification sent", "tasks", len(tasks), "email", user.Email)
			}
		}
	}

	// SMS notification — claim once per user per day
	if s.smsNotifier != nil && user.NotifySMS && user.Phone != "" {
		notif := entities.NewBatchNotification(user.ID, "sms", dueDate)
		claimed, err := s.notifRepo.Claim(ctx, notif)
		if err != nil {
			slog.Error("SMS claim error", "userID", user.ID, "error", err)
		} else if claimed {
			if err := s.smsNotifier.Send(ctx, user.Phone, subject, body); err != nil {
				slog.Error("SMS send error", "userID", user.ID, "phone", user.Phone, "error", err)
				if delErr := s.notifRepo.Delete(ctx, notif.ID); delErr != nil {
					slog.Error("Error releasing SMS claim", "error", delErr)
				}
			} else {
				slog.Info("SMS notification sent", "tasks", len(tasks), "phone", user.Phone)
			}
		}
	}

	// Push notification — claim once per user per day
	if s.pushNotifier != nil && user.NotifyPush {
		notif := entities.NewBatchNotification(user.ID, "push", dueDate)
		claimed, err := s.notifRepo.Claim(ctx, notif)
		if err != nil {
			slog.Error("Push claim error", "userID", user.ID, "error", err)
		} else if claimed {
			if err := s.pushNotifier.SendToUser(ctx, user.ID, subject, body); err != nil {
				slog.Error("Push send error", "userID", user.ID, "error", err)
				if delErr := s.notifRepo.Delete(ctx, notif.ID); delErr != nil {
					slog.Error("Error releasing push claim", "error", delErr)
				}
			} else {
				slog.Info("Push notification sent", "tasks", len(tasks), "userID", user.ID)
			}
		}
	}
}

func formatBatchBody(tasks []entities.Task) string {
	if len(tasks) == 1 {
		t := tasks[0]
		body := fmt.Sprintf("Your pool maintenance task \"%s\" is due today (%s).", t.Name, t.DueDate.Format("Jan 2, 2006"))
		if t.Description != "" {
			body += fmt.Sprintf("\n\nDetails: %s", t.Description)
		}
		return body
	}

	body := fmt.Sprintf("You have %d pool maintenance tasks due today:\n", len(tasks))
	for i, t := range tasks {
		body += fmt.Sprintf("\n%d. %s", i+1, t.Name)
		if t.Description != "" {
			body += fmt.Sprintf(" — %s", t.Description)
		}
	}
	return body
}
