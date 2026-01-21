export interface Notification {
  id: string;
  message: string;
  timestamp: number;
  type?: 'info' | 'success' | 'error' | 'warning';
}

export const mapEventTypeToType = (eventType: string): Notification['type'] => {
  switch (eventType.toLowerCase()) {
    case 'heart beat':
      return 'info';
    case 'breach':
    case 'phishing':
      return 'error';
    case 'alert':
      return 'warning';
    default:
      return 'success';
  }
};