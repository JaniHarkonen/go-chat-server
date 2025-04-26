import { gotoView } from "../index.js";
import ChatView from "./chatView.js";

  // Login screen component
export default function LoginView(props) {

    // Join the chatroom on login
  function onLogin(e) {
    e.preventDefault();
    gotoView(ChatView({
      id: "chat-view",
      username: Object.values(e.target)[0].value
    }));
  }

  const html = () => {
    return (
`
  <div id="${props.id}" class="w-100 h-100 d-flex d-flex-justify-content-center d-flex-align-items-center">
    <div class="w-50 h-50">
      <h1>Login</h1>
      <form onsubmit="onLogin(event)">
        <label class="w-100" for="login-username-input">Username:</label>
        <input class="w-100" id="login-username-input" type="text"></input>
        <!--<label class="w-100" for="login-password-input">Password:</label>
        <input class="w-100" id="login-password-input" type="password"></input>-->
      </form>
    </div>
  </div>
`
    );
  };

  return {
    html,
    onMount: () => {},
    scripts: [onLogin]
  };
}

