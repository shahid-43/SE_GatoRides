import React, { createContext, useState } from 'react';

const RideContext = createContext();

export const RideProvider = ({ children }) => {
  const [ridePayload, setRidePayload] = useState(null);

  return (
    <RideContext.Provider value={{ ridePayload, setRidePayload }}>
      {children}
    </RideContext.Provider>
  );
};

export default RideContext;