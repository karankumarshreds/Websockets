import React from 'react';
import { newSocketConnection } from 'config/socket';

const INPUT_STYLE = { padding: '10px', fontSize: 16, marginTop: 10 };

const Landing = () => {
  const [userDetails, setUserDetails] = React.useState({
    username: '',
    userId: '',
  });

  const submitHandler = (e) => {
    e.preventDefault();
    if (!userDetails.userId || !userDetails.username) return;
    const { webSocketConnection, error } = newSocketConnection(userDetails.userId);
    if (error) alert(error);
    else console.log({ webSocketConnection });
  };

  return (
    <div style={{ width: 600, margin: 'auto' }}>
      <h1>Enter your details</h1>
      <form onSubmit={submitHandler} style={{ display: 'flex', flexDirection: 'column' }}>
        <input
          placeholder="username"
          value={userDetails.username}
          onChange={(e) => setUserDetails((state) => ({ ...state, username: e.target.value }))}
          style={INPUT_STYLE}
        />
        <input
          placeholder="userid"
          value={userDetails.userId}
          onChange={(e) => setUserDetails((state) => ({ ...state, userId: e.target.value }))}
          style={INPUT_STYLE}
        />
        <button type={'submit'} style={INPUT_STYLE}>
          SUBMIT
        </button>
      </form>
    </div>
  );
};

export default Landing;
