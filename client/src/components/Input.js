import React from 'react';

const Input = ({ value, onChange, placeholder }) => {
  return <input className="form-control my-2" placeholder={placeholder} value={value} onChange={onChange} />;
};

export default Input;
