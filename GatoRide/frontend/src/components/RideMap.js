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

    // Updated fetchLocationSuggestions function to use Photon API
    const fetchLocationSuggestions = async (query, setSuggestions) => {
        if (!query) {
            setSuggestions([]);
            return;
        }
    
        try {
            const response = await fetch(
                `https://photon.komoot.io/api/?q=${encodeURIComponent(query)}`,
                {
                    headers: {
                        'Accept': 'application/json',
                    }
                }
            );
            const data = await response.json();
    
            // Check if data contains features
            if (data && Array.isArray(data.features) && data.features.length > 0) {
                setSuggestions(data.features.map((feature) => {
                    const properties = feature.properties;
                    const fullAddress = [
                        properties.name,
                        properties.street,
                        properties.postcode,
                        properties.city,
                        properties.country
                    ]
                        .filter(Boolean) // Remove undefined or null values
                        .join(', '); // Join with commas
    
                    return {
                        lat: feature.geometry.coordinates[1], // Latitude
                        lon: feature.geometry.coordinates[0], // Longitude
                        display_name: fullAddress || 'Unknown Location'
                    };
                }));
            } else {
                setSuggestions([]); // Clear suggestions if no results
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
        <div className="ride-map-container" >
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

            <div className="map-container" data-testid="map-container">
                <MapContainer center={[29.6516, -82.3248]} zoom={13} style={{ height: '400px', width: '100%' }}>
                    <TileLayer
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    />
                    {fromLocation && parseFloat(fromLocation.lat) && parseFloat(fromLocation.lon) && (
                        <Marker position={[parseFloat(fromLocation.lat), parseFloat(fromLocation.lon)]}>
                            <Popup>{fromLocation.display_name}</Popup>
                        </Marker>
                    )}
                    {toLocation && parseFloat(toLocation.lat) && parseFloat(toLocation.lon) && (
                        <Marker position={[parseFloat(toLocation.lat), parseFloat(toLocation.lon)]}>
                            <Popup>{toLocation.display_name}</Popup>
                        </Marker>
                    )}
                </MapContainer>
            </div>
        </div>
    );
};

export default RideMap;
