import React, { useState, useContext } from 'react';
import AuthContext from '../context/AuthContext';
import '../styles.css';  // Import the global styles
const SignupForm = () => {
  const { handleSignup } = useContext(AuthContext);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    username: '',
    password: '',
    location: '',
    latitude: null,
    longitude: null
  });
  const [suggestions, setSuggestions] = useState([]);
  const [showDropdown, setShowDropdown] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });

    if (name === 'location' && value.length > 2) {
      // Fetch location suggestions from OpenStreetMap
      fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${value}`)
        .then(response => response.json())
        .then(data => {
          setSuggestions(data);
          setShowDropdown(data.length > 0);
        });
    } else if (name === 'location') {
      setSuggestions([]);
      setShowDropdown(false);
    }
  };

  const handleLocationSelect = (location) => {
    const latitude = parseFloat(location.lat);
    const longitude = parseFloat(location.lon);
  
    setFormData({
      ...formData,
      location: location.display_name,
      latitude: latitude,
      longitude: longitude
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
        address: formData.location,  // Full address as a string
        latitude: formData.latitude,  // Numeric latitude
        longitude: formData.longitude // Numeric longitude
      }
    };
  
    console.warn("Sending Payload:", JSON.stringify(payload, null, 2)); // Debugging step
  
    try {
      await handleSignup(payload);
      alert('Sign up successful! Please check your email for verification.');
    } catch (error) {
      console.error("Signup Error:", error);
      alert('Error during signup');
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
