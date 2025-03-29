import React, { useState, useContext } from 'react';
import AuthContext from '../context/AuthContext';
import '../styles.css';  // Import the global styles
import axios from 'axios';
const SignupForm = () => {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    username: '',
    password: '',
    latitude: null,
    longitude: null,
    location: ''
  });
  const [suggestions, setSuggestions] = useState([]);
  const [showDropdown, setShowDropdown] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  
    if (name === 'location' && value.length > 2) {
      // Fetch location suggestions from PositionStack
      fetch(`http://api.positionstack.com/v1/forward?access_key=3c15e1cc65d40b3afb7635dadb08767a&query=${value}`)
        .then(response => response.json())
        .then(data => {
          if (data.data) {
            setSuggestions(data.data);
            setShowDropdown(data.data.length > 0);
          }
        })
        .catch(error => console.error("Error fetching location data:", error));
    } else if (name === 'location') {
      setSuggestions([]);
      setShowDropdown(false);
    }
  };
  

  const handleLocationSelect = (location) => {
    const latitude = parseFloat(location.latitude);
    const longitude = parseFloat(location.longitude);
  
    setFormData({
      ...formData,
      latitude: latitude,
      longitude: longitude,
      location: location.label  // PositionStack provides a "label" field for the full address
    });
  
    setSuggestions([]);
    setShowDropdown(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const payload = {
        name: formData.name,
        email: formData.email,
        username: formData.username,
        password: formData.password,
        location: {
            address: formData.location,
            latitude: formData.latitude ? parseFloat(formData.latitude) : 0,
            longitude: formData.longitude ? parseFloat(formData.longitude) : 0
        }
    };

    console.log("🔹 Sending Payload:", JSON.stringify(payload, null, 2));

    try {
        const response = await axios.post("http://localhost:5001/signup", payload, {
            headers: { "Content-Type": "application/json" }
        });

        console.log("✅ Signup Response:", response.data);
        alert("Sign up successful!");
    } catch (error) {
        console.error("❌ Signup Error:", error.response?.data || error);
        alert(`Error during signup: ${error.response?.data?.error || "Unknown error"}`);
    }
};

  
  

  return (
    <form onSubmit={handleSubmit} className="signup-form">
      <h2>Join GatoRides</h2>
      <p>Create your account to start ride-sharing.</p>
      <input type="text" name="name" placeholder="Name" onChange={handleChange} required />
      <input type="email" name="email" placeholder="Email" onChange={handleChange} required />
      <input type="text" name="username" placeholder="Username" onChange={handleChange} required />
      <input type="password" name="password" placeholder="Password" onChange={handleChange} required />
      
      <div className="dropdown">
        <input
          type="text"
          name="location"
          placeholder="Enter your location"
          value={formData.location}
          onChange={handleChange}
          onFocus={() => setShowDropdown(suggestions.length > 0)}
          required
        />
        {showDropdown && (
          <div className="dropdown-menu">
            {suggestions.map((suggestion, index) => (
              <div
                key={index}
                className="dropdown-item"
                onClick={() => handleLocationSelect(suggestion)}
              >
                {suggestion.display_name}
              </div>
            ))}
          </div>
        )}
      </div>
      
      <button type="submit">Sign Up</button>
    </form>
  );
};


export default SignupForm;
