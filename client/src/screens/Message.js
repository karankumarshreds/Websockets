import React from 'react';
import { useParams } from 'react-router-dom';
import { Socket } from 'config/socket';
// components
import OnlineUsers from 'components/OnlineUsers';
import Input from 'components/Input';

const Message = () => {
  const { userid, username } = useParams();
  const [message, setMessage] = React.useState('');
  const [selectedUser, setSelectedUser] = React.useState(null);
  const [socket, setSocket] = React.useState(null);

  React.useEffect(() => {
    if (!userid || !username) return;
    const _socket = new Socket(userid, username);
    setSocket(_socket);
    if (_socket.error) return alert(_socket.error);
    else console.log({ web_SocketConnection: _socket.web_SocketConnection });
    _socket.listen();
  }, []);

  const sendMessage = () => {
    socket.sendDirectMessage({
      sender: userid,
      userId: selectedUser,
      message,
    });
  };

  return (
    <div className="container" style={{ width: 600, margin: 'auto' }}>
      <h1>Message Screen</h1>
      <Input placeholder="Enter your message" value={message} onChange={(e) => setMessage(e.target.value)} />
      <div class="d-grid">
        <button className="btn btn-primary" type="button" onClick={sendMessage}>
          Send Message
        </button>
      </div>
      <OnlineUsers selectedUser={selectedUser} setSelectedUser={setSelectedUser} />
    </div>
  );
};

export default Message;
