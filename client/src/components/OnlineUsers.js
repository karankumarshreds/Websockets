import { EVENT_NAMES } from 'config/constants';
import React from 'react';

const selectedUserStyle = { fontWeight: 800, fontStyle: 'italic' };

const OnlineUsers = ({ selectedUser, setSelectedUser }) => {
  const [users, setUsers] = React.useState([]);

  React.useEffect(() => {
    eventEmitter.on(EVENT_NAMES.NEW_USER, (payload) => {
      setUsers(payload);
    });
  }, []);

  return (
    <div className="container mt-5">
      <h2>List of online users</h2>
      <ul>
        {users.map((each) => (
          <li
            key={each.userId}
            onClick={() => setSelectedUser(each.userId)}
            style={each.userId === selectedUser ? selectedUserStyle : { cursor: 'pointer' }}>
            {each.username}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default OnlineUsers;
