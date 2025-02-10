import React, { createContext, useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { login, signup, verifyEmail } from '../services/AuthService';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUser = () => {
      const storedUser = localStorage.getItem('user');
      if (storedUser) {
        setUser(JSON.parse(storedUser));
      }
      setLoading(false);
    };

    fetchUser();
  }, []);

  const handleLogin = async (email, password) => {
    try {
      const response = await login(email, password);
      setUser(response.user);
      localStorage.setItem('user', JSON.stringify(response.user));
      navigate('/dashboard');
    } catch (error) {
      console.error(error);
    }
  };

  const handleSignup = async (name, email, username, password) => {
    try {
      await signup(name, email, username, password);
      navigate('/login');
    } catch (error) {
      console.error(error);
    }
  };

  const handleVerifyEmail = async (token) => {
    try {
      await verifyEmail(token);
      setUser({ ...user, verified: true });
      localStorage.setItem('user', JSON.stringify({ ...user, verified: true }));
      navigate('/login');
    } catch (error) {
      console.error(error);
    }
  };

  const handleLogout = () => {
    setUser(null);
    localStorage.removeItem('user');
    navigate('/');
  };

  return (
    <AuthContext.Provider value={{
      user,
      loading,
      handleLogin,
      handleSignup,
      handleVerifyEmail,
      handleLogout
    }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
