package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type PushHandler struct {
	pushRepo repositories.PushSubscriptionRepository
	vapidPub string
}

func NewPushHandler(pushRepo repositories.PushSubscriptionRepository, vapidPub string) *PushHandler {
	return &PushHandler{pushRepo: pushRepo, vapidPub: vapidPub}
}

type pushSubscribeRequest struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

func (h *PushHandler) VAPIDPublicKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"publicKey": h.vapidPub})
}

func (h *PushHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID, err := services.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req pushSubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Endpoint == "" || req.Keys.P256dh == "" || req.Keys.Auth == "" {
		http.Error(w, "missing subscription fields", http.StatusBadRequest)
		return
	}

	sub := entities.NewPushSubscription(userID, req.Endpoint, req.Keys.P256dh, req.Keys.Auth)
	if err := h.pushRepo.Save(r.Context(), sub); err != nil {
		slog.Error("Error saving push subscription", "error", err)
		http.Error(w, "failed to save subscription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "subscribed"})
}

func (h *PushHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	userID, err := services.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Endpoint string `json:"endpoint"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.pushRepo.DeleteByEndpoint(r.Context(), userID, req.Endpoint); err != nil {
		slog.Error("Error deleting push subscription", "error", err)
		http.Error(w, "failed to remove subscription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "unsubscribed"})
}
