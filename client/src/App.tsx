import React, { useEffect, useState } from 'react';

import { Message } from './components/Message';
import { Login } from './components/Login';
import './App.css';

interface StateType {
  init: boolean,
  participants: Array<{}>,
  messages: Array<any>,
  name: string,
}

let socket: any = null;
const dev = process.env.NODE_ENV === 'development'

const App = () => {
  const [currentState, setCurrentState] = useState<StateType>({ init: true, participants: [], messages: [], name: "" });
  const [msgValue, setMsgValue] = useState('');

  const onParticipantSubmit = (e: any) =>  {
    e.preventDefault();
    const inputEl: HTMLInputElement = document.getElementById("formInput") as HTMLInputElement;
    const input = inputEl.value;

    if (input === "") return;

    const location = dev ? 'localhost:8081' : document.location.host;
    const url = `${document.location.protocol.replace("http", "ws")}//${location}/ws?name=${input}`;
    // init client websocket
    socket = new WebSocket(url);
    socket.onopen = (event: any) => {
      console.log('Opened socket');
    };
    socket.onerror = (event: any) => {
      console.log(event);
    }
    setCurrentState(prevState => ({ ...prevState, init: false, name: input, participants: [] }));
  }

  const onTextSubmit = (e: any) =>  {
    e.preventDefault();
    const inputEl: HTMLInputElement = document.getElementById("text-entry-input") as HTMLInputElement;
    let input = inputEl.value;

    if (input === "") return;

    // send message request to server
    socket.send(JSON.stringify({
      text: input,
      name: currentState.name,
      timestamp: new Date().toLocaleTimeString(),
    }));

    setMsgValue('');
  }

  const handleOnChange = (e: any) => {
    const targetValue = e.target.value;
    setMsgValue(targetValue);
  }

  useEffect(() => {
    if (!currentState.init) {
      // keep message box scrolled appropriately with new messages
      const mainElement = document.getElementById("conversation-main")!;
      mainElement.scrollTop = mainElement.scrollHeight;

      socket.onmessage = (event: any) => {
        const data = JSON.parse(event.data);

        if (!data.text) { // client added
          setCurrentState(prevState => ({ ...prevState, participants: data }));
        } else { // message added
          const { text, name, timestamp } = data;
          setCurrentState(prevState => ({ ...prevState, messages: [ ...prevState.messages, { text, name, timestamp }], }));
        }
      }
    }
  }, [currentState.init]);

  return (
    <div id="main">
      <div id="title">
        go-chat
      </div>
      {currentState.init &&
        <div id="box-init">
          <Login participantSubmit={onParticipantSubmit} />
        </div>
      }
      {!currentState.init &&
        <div id="box-main">
          <div style={{width: "100%", height: "100%"}}>
            <div id="top">
              <div id="participants-main">
                <b>{`Participants (${currentState.participants.length}):`}</b>
                {currentState.participants.map(participant =>
                  <div>{participant}</div>
                )}
              </div>
              <div id="conversation-main">
                {currentState.messages.map(message =>
                  <Message message={message} />
                )}
              </div>
            </div>
            <div id="bottom">
              <div id="text-entry-main">
                <form onSubmit={onTextSubmit} id="text-entry-form">
                  <input type="text" id="text-entry-input" placeholder="Type a message..." value={msgValue} onChange={handleOnChange} />
                </form>
              </div>
            </div>
          </div>
        </div>
      }
    </div>
  );
}

export default App;