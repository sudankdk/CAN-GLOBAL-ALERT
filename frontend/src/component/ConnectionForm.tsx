import React, { useState } from 'react';

interface ConnectionFormProps {
  onConnect: (id: string, email: string) => void;
}

export default function ConnectionForm({ onConnect }: ConnectionFormProps) {
  const [id, setId] = useState('');
  const [email, setEmail] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (id && email) {
      onConnect(id, email);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="p-4 bg-white rounded-lg shadow-md max-w-md mx-auto">
      <h2 className="text-lg font-semibold mb-4">Connect to Notifications</h2>
      <div className="mb-4">
        <label className="block text-sm font-medium mb-1">Notification ID</label>
        <input
          type="text"
          value={id}
          onChange={(e) => setId(e.target.value)}
          className="w-full p-2 border rounded"
          placeholder="e.g., test123"
          required
        />
      </div>
      <div className="mb-4">
        <label className="block text-sm font-medium mb-1">Email</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full p-2 border rounded"
          placeholder="e.g., user@example.com"
          required
        />
      </div>
      <button type="submit" className="w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600">
        Connect
      </button>
    </form>
  );
}