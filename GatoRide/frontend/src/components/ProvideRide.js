import React, { useState, useContext, useEffect } from 'react';
import AuthContext from '../context/AuthContext';
import '../styles.css';

const RideRequest = () => {
  const { user } = useContext(AuthContext); // Fetch user data from context

  // Initialize ride details
  const [rideDetails, setRideDetails] = useState({
    pickup: { address: '' }, // Initially empty
    dropoff: { address: '' },
    price: '',
  });

  const [suggestions, setSuggestions] = useState([]);
  const [showDropdown, setShowDropdown] = useState(false);

  // Set the pickup location from user data when the component mounts
  useEffect(() => {
    if (user?.location?.address) {
      setRideDetails((prevDetails) => ({
        ...prevDetails,
        pickup: { address: user.location.address },
      }));
    }
  }, [user]);

  // Fetch location suggestions
  const fetchLocationSuggestions = async (query) => {
    if (!query) {
      setSuggestions([]);
      return;
    }

    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}`
      );
      const data = await response.json();

      setSuggestions(data.slice(0, 5).map((item) => item.display_name));
      setShowDropdown(true);
    } catch (error) {
      console.error('Error fetching location suggestions:', error);
    }
  };

  // Handle dropoff location input change
  const handleLocationChange = (e) => {
    const { value } = e.target;
    setRideDetails((prevDetails) => ({
      ...prevDetails,
      dropoff: { address: value },
    }));

    fetchLocationSuggestions(value);
  };

  // Handle selecting a dropoff location from suggestions
  const handleLocationSelect = (selectedAddress) => {
    setRideDetails((prevDetails) => ({
      ...prevDetails,
      dropoff: { address: selectedAddress },
    }));
    setSuggestions([]);
    setShowDropdown(false);
  };

  // Handle price input change
  const handleChange = (e) => {
    setRideDetails((prevDetails) => ({
      ...prevDetails,
      price: e.target.value,
    }));
  };

  // Submit ride request
  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!rideDetails.dropoff.address) {
      alert('Please select a valid dropoff location.');
      return;
    }

    console.log('Submitting Ride Request:', rideDetails);

    // TODO: Send `rideDetails` to backend via API
    alert('Ride request submitted successfully!');
  };

  return (
    <div className="ride-request-container">
      <form onSubmit={handleSubmit}>
        <h2>Request a Ride</h2>

        <h3>Pickup Location</h3>
        <input type="text" name="pickup" value={rideDetails.pickup.address} disabled />

        <h3>Dropoff Location</h3>
        <div className="dropdown">
          <input
            type="text"
            name="dropoff"
            placeholder="Enter dropoff address"
            value={rideDetails.dropoff.address}
            onChange={handleLocationChange}
            onFocus={() => setShowDropdown(suggestions.length > 0)}
            required
          />
          {showDropdown && suggestions.length > 0 && (
            <div className="dropdown-menu">
              {suggestions.map((suggestion, index) => (
                <div
                  key={index}
                  className="dropdown-item"
                  onClick={() => handleLocationSelect(suggestion)}
                >
                  {suggestion}
                </div>
              ))}
            </div>
          )}
        </div>

        <h3>Price</h3>
        <input type="number" name="price" placeholder="Enter price" value={rideDetails.price} onChange={handleChange} required />

        <button type="submit">Submit Ride Request</button>
      </form>
    </div>
  );
};

export default RideRequest;
