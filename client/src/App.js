import React from 'react';
import { Routes, Route } from 'react-router-dom';
// screens
import Landing from 'screens/Landing';

const App = () => {
  return (
    <Routes>
      <Route path="/" element={<Landing />} />
    </Routes>
  );
};

export default App;
