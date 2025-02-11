import React from 'react';
import { Link } from 'react-router-dom';

const Home = () => {
  return (
    <div className="page-container">
      <div className="home-content">
        <h1>Welcome to GatoRides</h1>
        <p>Experience easy and reliable ride-sharing!</p>
        <div className="home-buttons">
          <Link to="/signup" className="btn">Sign Up</Link>
          <Link to="/login" className="btn">Login</Link>
          <Link to="/dashboard" className="btn">Dashboard</Link>
        </div>
      </div>
    </div>
  );
};

export default Home;
