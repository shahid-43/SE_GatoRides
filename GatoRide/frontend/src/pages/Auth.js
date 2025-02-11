const API_URL = "http://localhost:5000/api/auth";

const signup = async (userData) => {
  const response = await fetch(`${API_URL}/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(userData),
  });
  return response.json();
};

const login = async (credentials) => {
  const response = await fetch(`${API_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(credentials),
  });
  return response.json();
};

const getCurrentUser = async () => {
  const response = await fetch(`${API_URL}/current-user`, { credentials: "include" });
  return response.json();
};

export default { signup, login, getCurrentUser };