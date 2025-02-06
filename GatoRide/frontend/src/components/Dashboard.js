import React, { useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthContext from '../context/AuthContext';

const Dashboard = () => {
  const { user, handleLogout } = useContext(AuthContext);
  const navigate = useNavigate();

  const logoutHandler = () => {
    handleLogout();
    navigate('/');
  };

  return (
    <div>
      <h1>Welcome to your Dashboard</h1>
      {user ? (
        <div>
          <h2>User Details</h2>
          <p>Name: {user.name}</p>
          <p>Email: {user.email}</p>
          <p>Username: {user.username}</p>
          <p>Status: {user.verified ? 'Verified' : 'Not Verified'}</p>
        </div>
      ) : (
        <p>Loading user information...</p>
      )}
      <button onClick={logoutHandler}>Logout</button>
    </div>
  );
};

export default Dashboard;
