package services

import (
	"context"
	"log"
	"time"

	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type DemoCleanupService struct {
	userRepo      repositories.UserRepository
	sessionRepo   repositories.SessionRepository
	chemLogRepo   repositories.ChemistryLogRepository
	taskRepo      repositories.TaskRepository
	equipRepo     repositories.EquipmentRepository
	srRepo        repositories.ServiceRecordRepository
	chemRepo      repositories.ChemicalRepository
	taskNotifRepo repositories.TaskNotificationRepository
	interval      time.Duration
}

func NewDemoCleanupService(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	chemLogRepo repositories.ChemistryLogRepository,
	taskRepo repositories.TaskRepository,
	equipRepo repositories.EquipmentRepository,
	srRepo repositories.ServiceRecordRepository,
	chemRepo repositories.ChemicalRepository,
	taskNotifRepo repositories.TaskNotificationRepository,
	interval time.Duration,
) *DemoCleanupService {
	return &DemoCleanupService{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		chemLogRepo:   chemLogRepo,
		taskRepo:      taskRepo,
		equipRepo:     equipRepo,
		srRepo:        srRepo,
		chemRepo:      chemRepo,
		taskNotifRepo: taskNotifRepo,
		interval:      interval,
	}
}

func (s *DemoCleanupService) Start(ctx context.Context) {
	log.Printf("Demo cleanup scheduler started (interval: %s)", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run immediately on start
	s.cleanup(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Demo cleanup scheduler stopped")
			return
		case <-ticker.C:
			s.cleanup(ctx)
		}
	}
}

func (s *DemoCleanupService) cleanup(ctx context.Context) {
	users, err := s.userRepo.FindExpiredDemo(ctx, time.Now())
	if err != nil {
		log.Printf("Demo cleanup error finding expired users: %v", err)
		return
	}

	for _, user := range users {
		log.Printf("Cleaning up expired demo user: %s (%s)", user.Email, user.ID)

		// Delete all user data explicitly (no FK CASCADE on entity tables)
		chems, _ := s.chemRepo.FindAll(ctx, user.ID)
		for _, c := range chems {
			_ = s.chemRepo.Delete(ctx, user.ID, c.ID)
		}

		equips, _ := s.equipRepo.FindAll(ctx, user.ID)
		for _, e := range equips {
			records, _ := s.srRepo.FindByEquipmentID(ctx, user.ID, e.ID)
			for _, r := range records {
				_ = s.srRepo.Delete(ctx, user.ID, r.ID)
			}
			_ = s.equipRepo.Delete(ctx, user.ID, e.ID)
		}

		tasks, _ := s.taskRepo.FindAll(ctx, user.ID)
		for _, t := range tasks {
			_ = s.taskRepo.Delete(ctx, user.ID, t.ID)
		}

		logs, _ := s.chemLogRepo.FindAll(ctx, user.ID)
		for _, l := range logs {
			_ = s.chemLogRepo.Delete(ctx, user.ID, l.ID)
		}

		// Sessions are deleted via FK CASCADE, but clean up explicitly too
		_ = s.sessionRepo.DeleteByUserID(ctx, user.ID)

		// Finally delete the user
		if err := s.userRepo.Delete(ctx, user.ID); err != nil {
			log.Printf("Demo cleanup: failed to delete user %s: %v", user.ID, err)
		} else {
			log.Printf("Demo cleanup: deleted expired demo user %s", user.Email)
		}
	}
}
