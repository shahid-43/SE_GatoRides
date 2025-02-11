import React from 'react';
import { Routes, Route } from 'react-router-dom';
import RouteConfig from './routes/routes';

const App = () => {
  return (
    <Routes>
      {RouteConfig.map(({ path, component: Component }, index) => (
        <Route key={index} path={path} element={<Component />} />
      ))}
    </Routes>
  );
};


export default App;
