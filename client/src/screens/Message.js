import React from 'react';
import { useParams } from 'react-router-dom';
import { Socket } from 'config/socket';
import Input from 'components/Input';

const Message = () => {
  const { userid, username } = useParams();

  React.useEffect(() => {
    if (!userid || !username) return;
    const socket = new Socket(userid, username);

    if (socket.error) return alert(socket.error);
    else console.log({ webSocketConnection: socket.webSocketConnection });
    socket.listen();
  }, []);

  return (
    <div className="container" style={{ width: 600, margin: 'auto' }}>
      <h1>Message Screen</h1>
      <Input placeholder="Enter your message" />
      <div class="d-grid">
        <button className="btn btn-primary" type="button">
          Send Message
        </button>
      </div>
    </div>
  );
};

export default Message;
