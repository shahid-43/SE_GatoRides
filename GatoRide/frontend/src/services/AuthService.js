import axios from 'axios';
const API_URL = 'http://localhost:5000/api/auth'; 
export const signup = async (name, email, username, password) => {
  const response = await axios.post(`${API_URL}/signup`, { name, email, username, password });
  return response.data;
};
export const login = async (email, password) => {
  const response = await axios.post(`${API_URL}/login`, { email, password });
  const { token, message } = response.data;
  if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    return { user: { email, token }, message };
  }
  throw new Error(message);
};
export const verifyEmail = async (token) => {
  const response = await axios.get(`${API_URL}/verify-email/${token}`);
  return response.data;
};
