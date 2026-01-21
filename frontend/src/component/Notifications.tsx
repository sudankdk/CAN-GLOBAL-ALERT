import React, { useState, useEffect } from 'react';
import NotificationToast from './NotificationToast';
import type { Notification } from '../Interface/Notification';
import { mapEventTypeToType } from '../Interface/Notification';

interface NotificationsProps {
  id: string;
  email: string;
}

export default function Notifications({ id, email }: NotificationsProps) {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [eventSource, setEventSource] = useState<EventSource | null>(null);

  useEffect(() => {
    if (!id || !email) return;

    const es = new EventSource(`http://localhost:8080/notifications/live?id=${encodeURIComponent(id)}&email=${encodeURIComponent(email)}`);
    setEventSource(es);

    es.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const newNotification: Notification = {
          id: data.timestamp.toString(),
          message: data.message,
          timestamp: data.timestamp,
          type: mapEventTypeToType(data.eventType),
        };
        setNotifications((prev) => [newNotification, ...prev].slice(0, 10)); // Keep last 10
      } catch (error) {
        console.error('Error parsing notification:', error);
      }
    };

    es.onerror = (error) => {
      console.error('EventSource error:', error);
    };

    return () => {
      es.close();
    };
  }, [id, email]);

  const removeNotification = (id: string) => {
    setNotifications((prev) => prev.filter((n) => n.id !== id));
  };

  return (
    <div className="fixed bottom-0 right-0 p-4 space-y-2">
      {notifications.map((notification) => (
        <NotificationToast
          key={notification.id}
          notification={notification}
          onClose={() => removeNotification(notification.id)}
        />
      ))}
    </div>
  );
}