import React from 'react';
import { useParams } from 'react-router-dom';
import { newSocketConnection } from 'config/socket';
import Input from 'components/Input';

const Message = () => {
  const { userid, username } = useParams();

  React.useEffect(() => {
    if (!userid || !username) return;
    const { webSocketConnection, error } = newSocketConnection(userid);
    if (error) alert(error);
    else console.log({ webSocketConnection });
  }, []);

  return (
    <div className="container">
      <h1>Message Screen</h1>
      <Input placeholder="Enter your message" />
    </div>
  );
};

export default Message;
