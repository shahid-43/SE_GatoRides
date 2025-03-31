import React, { useState, useContext } from 'react';
import { useNavigate } from 'react-router-dom'; // Import useNavigate
import AuthContext from '../context/AuthContext';
import '../styles.css';  // Import the global styles

// function LoginForm() {
//   const [formData, setFormData] = useState({
//     email: '',
//     password: '',
//   });

const Login = () => {
  const { handleLogin } = useContext(AuthContext); // Use handleLogin from context
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const navigate = useNavigate(); // Initialize useNavigate

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await handleLogin(formData.email, formData.password);
      alert('Login successful');
      navigate('/'); // Redirect to the base URL
    } catch (error) {
      alert(error.response?.data?.message || 'Error logging in');
    }
  };

// const handleSubmit = async (e) => {
  //   e.preventDefault();
  //   try {
  //     console.warn(formData);
  //     console.warn(handleSignup);
  //     await handleSignup(formData.name, formData.email, formData.username, formData.password);
  //     alert('Sign up successful! Please check your email for verification.');
  //   } catch (error) {
  //     alert('Error during signup');
  //   }
  // };

  return (
    <form onSubmit={handleSubmit} >
      <h2>Login to GatoRides</h2>
      <p>Access your ride-sharing account.</p>
      <input type="email" name="email" placeholder="Email" onChange={handleChange} required />
      <input type="password" name="password" placeholder="Password" onChange={handleChange} required />
      <button type="submit">Login</button>
    </form>
  );
}

export default Login;
