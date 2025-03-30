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
    navigate('/ride-request'); // Navigate to RideRequest page
  };

  return (
    <div className="dashboard-container">
      <div className="left-column">
        <div className="user-details">
          <h2>User Details</h2>
          <p><strong>Email:</strong> {user.email}</p>
          <p><strong>Username:</strong> {user.username}</p>
        </div>

        <div className="additional-element">
          <h3>Additional Element</h3>
          <p>Content for the additional element goes here.</p>
        </div>

        <div className="actions">
          {/* Create Ride Button */}
          <button onClick={handleCreateRide} className="btn btn-primary">Create Ride</button>

          {/* Logout Button */}
          <button onClick={logoutHandler} className="btn btn-secondary">Logout</button>
        </div>
      </div>

      <div className="right-column">
        <RideMap />
      </div>
    </div>
  );
};

export default Dashboard;
