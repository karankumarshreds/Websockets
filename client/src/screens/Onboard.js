import React from 'react';
import { useNavigate } from 'react-router-dom';
import Input from 'components/Input';

const Landing = (props) => {
  const [userDetails, setUserDetails] = React.useState({
    username: '',
    userId: '',
  });
  const navigate = useNavigate();

  const submitHandler = (e) => {
    e.preventDefault();
    if (userDetails.userId && userDetails.username) {
      navigate(`/${userDetails.userId}/${userDetails.username}`);
    }
  };

  return (
    <div style={{ width: 600, margin: 'auto', marginTop: 100 }} className="container">
      <h1>Enter your details</h1>
      <form onSubmit={submitHandler} style={{ display: 'flex', flexDirection: 'column' }}>
        <Input
          className="form-control my-2"
          placeholder="Username"
          value={userDetails.username}
          onChange={(e) => setUserDetails((state) => ({ ...state, username: e.target.value }))}
        />
        <Input
          placeholder="Userid"
          className="form-control my-2"
          value={userDetails.userId}
          onChange={(e) => setUserDetails((state) => ({ ...state, userId: e.target.value }))}
        />
        <button type={'submit'} className="btn btn-primary">
          SUBMIT
        </button>
      </form>
    </div>
  );
};

export default Landing;
