package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sudankdk/cga/internal/domain"
	my_redis "github.com/sudankdk/cga/internal/redis"
	"github.com/sudankdk/cga/internal/service"
)

type Handler struct {
	services service.NotificationServce
}

func NewHandler(svc service.NotificationServce) *Handler {
	return &Handler{
		services: svc,
	}
}

func (h *Handler) HandleLiveNotification(w http.ResponseWriter, r *http.Request) {
	// Query params
	notificationID := r.URL.Query().Get("id")
	email := r.URL.Query().Get("email")

	if notificationID == "" || email == "" {
		http.Error(w, "missing notification id or email", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	clientChan := make(chan domain.Notification, 101)

	h.services.Register(notificationID, email, clientChan)
	defer func() {
		h.services.Remove(notificationID, email)
		close(clientChan)
	}()

	// Initial flush to establish stream
	flusher.Flush()

	for {
		select {
		case msg, ok := <-clientChan:
			if !ok {
				return
			}

			payload, err := json.Marshal(msg)
			if err != nil {
				continue
			}

			w.Write([]byte("data: "))
			w.Write(payload)
			w.Write([]byte("\n\n"))

			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

func (h *Handler) HandleBroadcastMessage(w http.ResponseWriter, r *http.Request) {
	notificationID := r.URL.Query().Get("id")
	if notificationID == "" {
		http.Error(w, "missing notification id", http.StatusBadRequest)
		return
	}

	var message domain.Notification
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&message); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if message.EventType == "" || message.Message == "" {
		http.Error(w, "event_type and message are required", http.StatusBadRequest)
		return
	}

	redisMsg := domain.RedisMessage{
		NotificationId: notificationID,
		EventType:      message.EventType,
		Message:        message.Message,
		Timestamp:      message.Timestamp,
	}

	data, err := json.Marshal(redisMsg)
	if err != nil {
		http.Error(w, "failed to marshal notification", http.StatusInternalServerError)
		return
	}

	if err := my_redis.RedisClient.Publish(context.Background(), "notifications", data).Err(); err != nil {
		http.Error(w, "failed to publish notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
