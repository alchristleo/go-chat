import React, { FC } from 'react';

import './styles/Login.css';

interface LoginProps {
  participantSubmit?(e: any): void;
}

export const Login: FC<LoginProps> = ({ participantSubmit }) => {
  return (
    <div id="loginWrapper">
      <form onSubmit={participantSubmit}>
        <input type="text"
          placeholder="Enter your name here"
          id="formInput" />
      </form>
    </div>
  );
};