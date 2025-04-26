import { getChatRoute } from "../CONFIG.js";

  // Login screen component
export default function ChatView(props) {
  let socket;

  function receiveMessage(e) {
    console.log(e.data);
  }

  function sendMessage(messge) {
    socket.send(message);
  }

  function connectToChatroom() {
    socket = new WebSocket(getChatRoute());
    socket.onmessage = receiveMessage;
  }

  const html = () => {
    return (
`
  <div id="${props.id}" class="w-100 h-100 d-flex d-flex-justify-content-center d-flex-align-items-center">
    <div class="w-50 h-50">
      <h1>Your username is ${props.username}</h1>
    </div>
  </div>
`
    );
  };

  return {
    html,
    onMount: () => {
      console.log("mounted");
      connectToChatroom();
    },
    scripts: []
  };
}

