const iHeadClientInfo = 1;     // Client is sending its info upon joining.
const iHeadNameChange = 2;     // Client is requesting a name change.
const iHeadChatInput = 3;      // Client is sending a chat message/command
const oHeadNameChange = 4;     // Server is confirming a name change.
const oHeadCompleteUpdate = 5; // Server is sending a complete update containing latest messages, active and inactive users.
const oHeadDeltaUpdate = 6;    // Server is sending the latest delta snapshot of messages, active and inactive users.

const encoder = new TextEncoder();

function arrayInsertAt(index, array, ...element) {
  return [...array.slice(0, index), ...element, ...array.slice(index)];
}

function createMessage(header) {
  return [header];
}

function writeString(string, array, offset = 0) {
  return arrayInsertAt(offset, array, ...encoder.encode(string), 0);
}

function finalizeMessage(array) {
  const arrayBuffer = new Uint8Array(array.length);
  arrayBuffer.set(array);
  return arrayBuffer;
}

  // Client to server: Send client username upon joining
export function sendClientInfo(socket, data) {
  let response = createMessage(iHeadClientInfo);
  response = writeString(data.username, response, 1);
  socket.send(finalizeMessage(response).buffer);
}
