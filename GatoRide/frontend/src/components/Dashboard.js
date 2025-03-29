import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthContext from '../context/AuthContext';
import RideMap from './RideMap';
import '../styles.css';  // Import global styles

const Dashboard = () => {
  const { user, handleLogout } = useContext(AuthContext);
  const navigate = useNavigate();

  const logoutHandler = () => {
    handleLogout();
    navigate('/');
  };

  const handleCreateRide = () => {
    navigate('/ride-request');  // Navigate to RideRequest page
  };

  return (
    <div className="container">
      <h1>GatoRides Dashboard</h1>
      <p>Manage your ride-sharing experience effortlessly.</p>

      {user ? (
        <div>
          <h2>User Details</h2>
          <p><strong>Email:</strong> {user.email}</p>
          <p><strong>Username:</strong> {user.username}</p>

          <h2>Plan Your Ride</h2>
          <RideMap />

          {/* Create Ride Button */}
          <button onClick={handleCreateRide} className="btn btn-primary">Create Ride</button>

          {/* Logout Button */}
          <button onClick={logoutHandler} className="btn btn-secondary">Logout</button>
        </div>
      ) : (
        <p>Loading user information...</p>
      )}
    </div>
  );
};

export default Dashboard;
