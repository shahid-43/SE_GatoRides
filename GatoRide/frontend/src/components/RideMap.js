import React, { useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';

// Fix for default marker icons in React-Leaflet
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
    iconUrl: require('leaflet/dist/images/marker-icon.png'),
    shadowUrl: require('leaflet/dist/images/marker-shadow.png')
});

const RideMap = () => {
    const [fromLocation, setFromLocation] = useState(null);
    const [toLocation, setToLocation] = useState(null);
    const [fromSuggestions, setFromSuggestions] = useState([]);
    const [toSuggestions, setToSuggestions] = useState([]);
    const [error, setError] = useState('');

    // Fetch location suggestions
    const fetchLocationSuggestions = async (query, setSuggestions) => {
        if (!query) {
            setSuggestions([]);
            return;
        }
    
        try {
            const response = await fetch(
                `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}&countrycodes=us&state=Florida&limit=5`,
                {
                    headers: {
                        'Accept': 'application/json',
                        'User-Agent': 'GatorRides_App/1.0'
                    }
                }
            );
            const data = await response.json();
    
            // Check if data is an array and contains elements
            if (Array.isArray(data) && data.length > 0) {
                setSuggestions(data.map((item) => ({
                    lat: parseFloat(item.lat),
                    lon: parseFloat(item.lon),
                    display_name: item.display_name
                })));
            } else {
                setSuggestions([]); // If no results, clear suggestions
            }
        } catch (error) {
            console.error('Error fetching location suggestions:', error);
            setSuggestions([]); // Clear suggestions in case of an error
        }
    };
    
    // Handle user typing in input fields
    const handleInputChange = (e, setSuggestions) => {
        fetchLocationSuggestions(e.target.value, setSuggestions);
    };

    // Handle selecting a location from the suggestions
    const handleSelectLocation = (selectedLocation, setLocation, setSuggestions, inputId) => {
        setLocation(selectedLocation);
        setSuggestions([]);
        document.getElementById(inputId).value = selectedLocation.display_name;
    };

    // Handle search submission
    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        if (!fromLocation || !toLocation) {
            setError('Please select valid locations.');
            return;
        }

        console.log('From:', fromLocation);
        console.log('To:', toLocation);
    };

    return (
        <div className="ride-map-container">
            <form onSubmit={handleSubmit} className="location-form">
                <div className="input-group">
                    <label htmlFor="from">From:</label>
                    <input
                        type="text"
                        id="from"
                        name="from"
                        placeholder="Enter pickup location"
                        onChange={(e) => handleInputChange(e, setFromSuggestions)}
                        required
                    />
                    {fromSuggestions.length > 0 && (
                        <div className="dropdown-menu">
                            {fromSuggestions.map((location, index) => (
                                <div
                                    key={index}
                                    className="dropdown-item"
                                    onClick={() => handleSelectLocation(location, setFromLocation, setFromSuggestions, 'from')}
                                >
                                    {location.display_name}
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                <div className="input-group">
                    <label htmlFor="to">To:</label>
                    <input
                        type="text"
                        id="to"
                        name="to"
                        placeholder="Enter destination"
                        onChange={(e) => handleInputChange(e, setToSuggestions)}
                        required
                    />
                    {toSuggestions.length > 0 && (
                        <div className="dropdown-menu">
                            {toSuggestions.map((location, index) => (
                                <div
                                    key={index}
                                    className="dropdown-item"
                                    onClick={() => handleSelectLocation(location, setToLocation, setToSuggestions, 'to')}
                                >
                                    {location.display_name}
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                <button type="submit" className="search-button">Search Route</button>
            </form>

            {error && <div className="error-message" role="alert">{error}</div>}

            <div className="map-container">
                <MapContainer center={[29.6516, -82.3248]} zoom={13} style={{ height: '400px', width: '100%' }}>
                    <TileLayer
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    />
                    {fromLocation && (
                        <Marker position={[fromLocation.lat, fromLocation.lon]}>
                            <Popup>{fromLocation.display_name}</Popup>
                        </Marker>
                    )}
                    {toLocation && (
                        <Marker position={[toLocation.lat, toLocation.lon]}>
                            <Popup>{toLocation.display_name}</Popup>
                        </Marker>
                    )}
                </MapContainer>
            </div>
        </div>
    );
};

export default RideMap;
