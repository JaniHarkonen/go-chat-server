import { getChatRoute } from "../config.js";
import { iHeadCompleteUpdate, readCompleteUpdate, sendChatInput, sendClientInfo } from "../net.js";

  // Login screen component
export default function ChatView(props) {
  let socket;
  let users = {};

  function receiveMessage(e) {
    e.data.arrayBuffer().then((buffer) => {
      const dataView = new DataView(buffer);
      const header = dataView.getUint8();

      switch( header ) {
        case iHeadCompleteUpdate: {
          const chatState = readCompleteUpdate(dataView, 1);
          users = chatState.users;
          console.log(dataView);
          console.log(chatState);

          for( let chatMessage of chatState.messages ) {
            const messageElement = document.createElement("div");
            messageElement.innerHTML = (`
              <b>${users[chatMessage.userId] ? users[chatMessage.userId] : "ERROR"}:</b> ${chatMessage.message}
            `);

            const chatElement = document.getElementById("chat-content-messages");
            chatElement.appendChild(messageElement);
          }
        } break;

        default: {
          console.log("Attempting to handle message with invalid header '" + header + "'!");
        } break;
      }
    });
  }

  function sendMessage(e) {
    e.preventDefault();
    const messageInput = Object.values(e.target)[0];
    sendChatInput(socket, {input: messageInput.value});
    messageInput.value = "";
  }

  function connectToChatroom() {
    socket = new WebSocket(getChatRoute());
    socket.onopen = () => {
      const caption = document.getElementById("chat-caption-h1");
      caption.innerHTML = "Connected to chatroom";
      sendClientInfo(socket, {username: props.username});
    };
    socket.onclose = (e) => {
      const caption = document.getElementById("chat-caption-h1");
      caption.innerHTML = "Connection lost! Reason: " + e.code;
    };
    socket.onmessage = receiveMessage;
  }

  const html = () => {
    return (`
      <div id="${props.id}" class="w-100 h-100 d-flex d-flex-justify-content-center d-flex-align-items-center">
        <div class="w-50 h-100">
          <div class="chat-content-grid h-100 w-100">
            <h1 id="chat-caption-h1">Connecting to <code>${getChatRoute()}</code> as "${props.username}"</h1>
            <div id="chat-content-messages" class="overflow-y-auto">

            </div>
            <div id="chat-content-input" class="padding-16px">
              <form onsubmit="sendMessage(event)" class="w-100 d-flex gap-8px">
                <input class="w-100" type="text"></input>
                <button>Send</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    `);
  };

  return {
    html,
    onMount: () => {
      connectToChatroom();
    },
    scripts: [sendMessage]
  };
}

