  
import React, { FC } from 'react';

import './styles/Message.css';

interface MessageProps {
  text: string,
  name: string,
  timestamp: string
}

interface PropsType {
  message: MessageProps
}

export const Message: FC<PropsType> = (props: PropsType) => {
  const { name, timestamp, text } = props.message;

  return (
    <div className="message">
      <div className="messageTop">
        <span className="clientName">{name}</span>
        <span className="timestamp">{timestamp}</span>
      </div>
      <div className="messageText">
        {text}
      </div>
    </div>
  );
}