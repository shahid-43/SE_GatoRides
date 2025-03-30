import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthContext from '../context/AuthContext';
import RideContext from '../context/RideContext';
import RideMap from './RideMap';
import '../styles.css';  // Import global styles

const Dashboard = () => {
  const { user, handleLogout } = useContext(AuthContext);
  const { ridePayload } = useContext(RideContext);
  const navigate = useNavigate();

  const logoutHandler = () => {
    handleLogout();
    navigate('/');
  };

  const handleCreateRide = () => {
    navigate('/ride-request'); // No need to pass setRidePayload
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
          <h3>Ride Details</h3>
          {ridePayload ? (
            <div>
              <p><strong>Pickup:</strong> {ridePayload.pickup.address}</p>
              <p><strong>Dropoff:</strong> {ridePayload.dropoff.address}</p>
              <p><strong>Price:</strong> ${ridePayload.price}</p>
              <p><strong>Date:</strong> {ridePayload.date}</p>
            </div>
          ) : (
            <p>No ride details available.</p>
          )}
        </div>

        <div className="actions">
          <button onClick={handleCreateRide} className="btn btn-primary">Create Ride</button>
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
