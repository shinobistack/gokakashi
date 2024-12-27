import { useState } from 'react';
import shinobiStackLogo from '../assets/shinobistack.jpeg';

const AdminLoginPage = () => {
  const [adminSecret, setAdminSecret] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    // Handle admin login logic here
    console.log('Admin Secret:', adminSecret);
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md">
        <img 
            src={shinobiStackLogo}
            alt="Admin Logo" 
            className="w-full h-auto mb-4" 
          />
        <h2 className="text-2xl font-bold text-center">gokakshi</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <input
              type="password"
              id="admin-secret"
              value={adminSecret}
              onChange={(e) => setAdminSecret(e.target.value)}
              required
              className="block w-full px-4 py-2 mt-1 border rounded-md focus:outline-none focus:ring focus:ring-blue-300"
              placeholder="Enter your admin secret"
            />
          </div>
          <button
            type="submit"
            className="w-full px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring focus:ring-blue-300"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  );
};

export default AdminLoginPage;
