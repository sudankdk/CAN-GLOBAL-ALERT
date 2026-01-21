package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sudankdk/cga/internal/domain"
	my_redis "github.com/sudankdk/cga/internal/redis"
)

type NotificationServce struct {
	mu      sync.Mutex
	clients map[string]map[string]chan domain.Notification // Map<NotificaionId, Map<Subscriber, Channel>>

}

func (n *NotificationServce) Register(id, email string, channel chan domain.Notification) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, ok := n.clients[id]; !ok {
		n.clients[id] = make(map[string]chan domain.Notification)
	}
	n.clients[id][email] = channel
}

func (n *NotificationServce) Remove(id, email string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if client, ok := n.clients[id]; ok {
		if channel, exist := client[email]; exist {
			close(channel)
			delete(client, email)
		}
		if len(client) == 0 {
			delete(n.clients, id)
		}
	}
}

func (n *NotificationServce) Broadcast(id string, message domain.Notification) error {
	n.mu.Lock()
	clients, exists := n.clients[id]

	n.mu.Unlock()
	if !exists {
		log.Printf("Notificaion not found: notificaionId=%s", id)
		return fmt.Errorf("notificaion %s not found", id)
	}
	for emails, ch := range clients {
		go func(email string, ch chan domain.Notification) {
			select {

			case ch <- message:
			default:
				log.Printf("Removing client due to failed broadcast: notificationId=%s, subscriberEmail=%s", id, email)
				n.Remove(id, email)
			}
		}(emails, ch)
	}
	log.Panicf("broadcasting messages: %d", id)
	return nil
}

func (n *NotificationServce) BroadcastToClients(id, email string, messsage domain.Notification) error {
	n.mu.Lock()
	clients, exist := n.clients[id]
	if !exist {
		log.Printf("Notificaion not found: notificaionId=%s", id)
		return fmt.Errorf("notificaion %s not found", id)
	}
	ch, exist := clients[email]
	n.mu.Unlock()
	if !exist {
		return fmt.Errorf("subscriber %s not found for notification %s", email, id)
	}
	go func(email string, ch chan domain.Notification) {
		select {
		case ch <- messsage:
		default:
			log.Printf("Removing client due to failed broadcast: notificationId=%s, subscriberEmail=%s", id, email)
			n.Remove(id, email)
		}
	}(email, ch)
	log.Printf("broadcasting message: notificationId=%s , subscriberId=%s", id, email)
	return nil
}

func (n *NotificationServce) TriggerEvents(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		n.mu.Lock()
		notificaionIds := make([]string, 0, len(n.clients))
		for id := range n.clients {
			notificaionIds = append(notificaionIds, id)
		}
		n.mu.Unlock()
		for _, notificaionId := range notificaionIds {
			n.Broadcast(notificaionId, domain.Notification{
				EventType: "Heart Beat",
				Message:   "Ping",
				Timestamp: time.Now().Unix(),
			})
		}
	}
}

func (n *NotificationServce) SubToChannel(channel string) {
	pubsub := my_redis.RedisClient.Subscribe(context.Background(), channel)
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var message domain.RedisMessage
		err := json.Unmarshal([]byte(msg.Payload), &message)
		if err != nil {
			log.Println("Error processing message:", err)
			continue
		}

		broadcastChanges := domain.Notification{
			EventType: message.EventType,
			Message:   message.Message,
			Timestamp: message.Timestamp,
		}
		n.Broadcast(message.NotificationId, broadcastChanges)
		if message.SubscriberEmail != "" {
			userSpecificUpdate := domain.Notification{
				EventType: message.EventType,
				Message:   message.Message,
				Timestamp: message.Timestamp,
			}
			n.BroadcastToClients(message.NotificationId, message.SubscriberEmail, userSpecificUpdate)
		}
	}

}
