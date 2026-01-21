import { useState, useEffect } from 'react';
import type { Notification } from '../Interface/Notification';

interface NotificationToastProps {
  notification: Notification;
  onClose: () => void;
}

export default function NotificationToast({ notification, onClose }: NotificationToastProps) {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    setIsVisible(true);
    const timer = setTimeout(() => {
      setIsVisible(false);
      onClose();
    }, 6000);

    return () => clearTimeout(timer);
  }, [notification.id, onClose]);

  if (!isVisible) {
    return null;
  }

  const getTypeStyles = () => {
    switch (notification.type) {
      case 'success':
        return 'bg-green-500 text-white shadow-lg shadow-green-500/50';
      case 'error':
        return 'bg-red-500 text-white shadow-lg shadow-red-500/50';
      case 'warning':
        return 'bg-yellow-500 text-white shadow-lg shadow-yellow-500/50';
      default:
        return 'bg-blue-500 text-white shadow-lg shadow-blue-500/50';
    }
  };

  return (
    <div
      className={`max-w-sm animate-in slide-in-from-bottom-5 duration-300 ${getTypeStyles()} rounded-lg p-4 pointer-events-auto`}
    >
      <div className="flex items-start gap-3">
        <div className="text-xl mt-0.5">
          {notification.type === 'success' && '✓'}
          {notification.type === 'error' && '✕'}
          {notification.type === 'warning' && '⚠'}
          {!notification.type && 'ℹ'}
        </div>
        <div className="flex-1">
          <p className="font-medium text-sm">{notification.type?.toUpperCase() || 'NOTIFICATION'}</p>
          <p className="text-sm mt-1 opacity-95">{notification.message}</p>
        </div>
        <button
          onClick={() => {
            setIsVisible(false);
            onClose();
          }}
          className="text-xl leading-none opacity-70 hover:opacity-100 transition-opacity flex-shrink-0"
          aria-label="Close"
        >
          ×
        </button>
      </div>
    </div>
  );
}
