import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import AuthContext from '../context/AuthContext';
import Dashboard from '../components/Dashboard';

const Home = () => {
  const { user } = useContext(AuthContext);

  return (
    <div className="page-container">
      {user ? (
        <Dashboard />
      ) : (
        <div className="home-content">
          <h1>Welcome to GatoRides</h1>
          <p>Experience easy and reliable ride-sharing!</p>
          <div className="auth-buttons">
            <Link to="/login" className="btn btn-primary">Login</Link>
            <Link to="/signup" className="btn btn-secondary">Sign Up</Link>
          </div>
        </div>
      )}
    </div>
  );
};

export default Home;
