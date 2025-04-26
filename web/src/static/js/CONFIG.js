export const CHAT_PROTOCOL = "http";
export const CHAT_ADDRESS = "localhost";
export const CHAT_PORT = "12345";
export const CHAT_ROUTE = "/chat";

export function getChatRoute() {
  return CHAT_PROTOCOL + "://" + CHAT_ADDRESS + ":" + CHAT_PORT + CHAT_ROUTE;
}
