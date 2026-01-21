package domain

type Notification struct {
	EventType string `json:"eventType"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type RedisMessage struct {
	NotificationId  string `json:"notitification_id"`
	SubscriberEmail string `json:"subscriberEmail,omitempty"`
	EventType       string `json:"eventType"`
	Message         string `json:"message"`
	Timestamp       int64  `json:"timestamp"`
}
