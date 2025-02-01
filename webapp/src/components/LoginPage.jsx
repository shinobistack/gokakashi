import { useState } from "react";
import PropTypes from "prop-types";
import shinobiStackLogo from "../assets/shinobistack.png";

const LoginPage = ({ login }) => {
  const [adminSecret, setAdminSecret] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    login({ adminSecret: adminSecret });
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-300">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-3xl shadow-md">
        <img
          src={shinobiStackLogo}
          alt="Admin Logo"
          className="w-full h-auto mb-4 rounded-3xl"
        />
        <h2 className="text-2xl font-bold font-mono text-center">gokakashi</h2>
        <h3 className="text-1xl font-thin text-center">
          The centralized security platform
        </h3>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <input
              type="password"
              id="admin-secret"
              value={adminSecret}
              onChange={(e) => setAdminSecret(e.target.value)}
              required
              className="block w-full px-4 py-2 mt-1 border rounded-3xl focus:outline-none focus:ring focus:ring-blue-300"
              placeholder="Enter your admin secret"
            />
          </div>
          <button
            type="submit"
            className="w-full px-4 py-2 text-white bg-blue-600 rounded-3xl hover:bg-blue-700 focus:outline-none focus:ring focus:ring-blue-300"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  );
};

LoginPage.propTypes = {
  login: PropTypes.func.isRequired,
};

export default LoginPage;
