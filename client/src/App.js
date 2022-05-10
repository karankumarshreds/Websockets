import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Message from 'screens/Message';
// screens
import Onboard from 'screens/Onboard';

const App = () => {
  return (
    <div style={{ marginTop: 100 }}>
      <Routes>
        <Route path="/" element={<Onboard />} />
        <Route path="/:userid/:username" element={<Message />} />
        <Route path="/mychats/:userid" element={<div>Hey!!!</div>} />
      </Routes>
    </div>
  );
};

export default App;
