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
      </Routes>
    </div>
  );
};

export default App;
