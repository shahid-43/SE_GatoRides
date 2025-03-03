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
    const [error, setError] = useState('');

    const searchLocation = async (query) => {
        try {
            const response = await fetch(
                `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}&countrycodes=us&state=Florida&limit=1`,
                {
                    headers: {
                        'Accept': 'application/json',
                        'User-Agent': 'GatorRides_App/1.0'
                    }
                }
            );
            const data = await response.json();
            
            if (data.length > 0) {
                return {
                    lat: parseFloat(data[0].lat),
                    lon: parseFloat(data[0].lon),
                    display_name: data[0].display_name
                };
            }
            throw new Error('Location not found');
        } catch (error) {
            throw new Error('Location not found. Please try again.');
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        
        try {
            const from = e.target.elements.from.value;
            const to = e.target.elements.to.value;

            if (!from || !to) {
                return;
            }

            const fromResult = await searchLocation(from);
            const toResult = await searchLocation(to);

            if (fromResult && toResult) {
                setFromLocation(fromResult);
                setToLocation(toResult);
            }
        } catch (error) {
            setError(error.message);
        }
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
                        required
                    />
                </div>
                <div className="input-group">
                    <label htmlFor="to">To:</label>
                    <input
                        type="text"
                        id="to"
                        name="to"
                        placeholder="Enter destination"
                        required
                    />
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