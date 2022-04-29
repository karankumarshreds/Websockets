import React from 'react';
import { useParams } from 'react-router-dom';
import { newSocketConnection } from 'config/socket';

const Message = () => {
  const { userid, username } = useParams();

  React.useEffect(() => {
    if (!userid || !username) return;
    const { webSocketConnection, error } = newSocketConnection(userid);
    if (error) alert(error);
    else console.log({ webSocketConnection });
  }, []);

  return (
    <div>
      <h1>Message Screen</h1>
    </div>
  );
};

export default Message;
