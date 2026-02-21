package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type WebPushNotifier struct {
	pushRepo   repositories.PushSubscriptionRepository
	vapidPub   string
	vapidPriv  string
	vapidEmail string
}

func NewWebPushNotifier(pushRepo repositories.PushSubscriptionRepository, vapidPub, vapidPriv, vapidEmail string) *WebPushNotifier {
	return &WebPushNotifier{
		pushRepo:   pushRepo,
		vapidPub:   vapidPub,
		vapidPriv:  vapidPriv,
		vapidEmail: vapidEmail,
	}
}

type pushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Tag   string `json:"tag"`
	URL   string `json:"url"`
}

// SendToUser sends a push notification to all subscriptions for the given user.
func (n *WebPushNotifier) SendToUser(ctx context.Context, userID uuid.UUID, subject, body string) error {
	subs, err := n.pushRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("finding push subscriptions: %w", err)
	}
	if len(subs) == 0 {
		return nil
	}

	payload, err := json.Marshal(pushPayload{
		Title: subject,
		Body:  body,
		Tag:   "poolvibes-tasks",
		URL:   "/",
	})
	if err != nil {
		return fmt.Errorf("marshaling push payload: %w", err)
	}

	var lastErr error
	for _, sub := range subs {
		if err := n.sendToSubscription(ctx, userID, sub, payload); err != nil {
			slog.Error("Push notification failed", "endpoint", sub.Endpoint, "error", err)
			lastErr = err
		}
	}
	return lastErr
}

func (n *WebPushNotifier) sendToSubscription(ctx context.Context, userID uuid.UUID, sub entities.PushSubscription, payload []byte) error {
	s := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}

	resp, err := webpush.SendNotification(payload, s, &webpush.Options{
		Subscriber:      n.vapidEmail,
		VAPIDPublicKey:  n.vapidPub,
		VAPIDPrivateKey: n.vapidPriv,
		TTL:             86400,
	})
	if err != nil {
		return fmt.Errorf("sending web push: %w", err)
	}
	defer resp.Body.Close()

	// 410 Gone or 404 means the subscription is no longer valid
	if resp.StatusCode == 410 || resp.StatusCode == 404 {
		slog.Info("Push subscription expired, removing", "endpoint", sub.Endpoint)
		_ = n.pushRepo.DeleteByEndpoint(ctx, userID, sub.Endpoint)
		return nil
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("push service returned status %d", resp.StatusCode)
	}

	return nil
}
