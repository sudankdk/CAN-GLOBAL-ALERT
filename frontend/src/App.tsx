import { useState } from 'react';
import ConnectionForm from './component/ConnectionForm';
import Notifications from './component/Notifications';

function App() {
  const [connected, setConnected] = useState(false);
  const [id, setId] = useState('');
  const [email, setEmail] = useState('');

  const handleConnect = (newId: string, newEmail: string) => {
    setId(newId);
    setEmail(newEmail);
    setConnected(true);
  };

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      {!connected ? (
        <ConnectionForm onConnect={handleConnect} />
      ) : (
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Connected to Notifications</h1>
          <p>ID: {id} | Email: {email}</p>
          <Notifications id={id} email={email} />
        </div>
      )}
    </div>
  );
}

export default App;
