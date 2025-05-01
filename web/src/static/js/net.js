export const iHeadNameChange = 4;     // Server is confirming a name change.
export const iHeadCompleteUpdate = 5; // Server is sending a complete update containing latest messages, active and inactive users.
export const iHeadDeltaUpdate = 6;    // Server is sending the latest delta snapshot of messages, active and inactive users.

const oHeadClientInfo = 1;     // Client is sending its info upon joining.
const oHeadNameChange = 2;     // Client is requesting a name change.
const oHeadChatInput = 3;      // Client is sending a chat message/command

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
  return arrayBuffer.buffer;
}

function readUint32(dataView, offset) {
  return {value: dataView.getUint32(offset), offset: offset + 4};
}

function readString(dataView, offset) {
  let string = "";
  let charAt;

  while( (charAt = dataView.getUint8(offset)) != 0 ) {
    string += String.fromCharCode(charAt);
    offset++;
  }

  offset++;

  return {value: string, offset};
}

function readUserId(dataView, offset) {
  return {value: dataView.getBigUint64(offset), offset: offset + 8};
}

  // Client to server: Send client username upon joining
export function sendClientInfo(socket, data) {
  let response = createMessage(oHeadClientInfo);
  response = writeString(data.username, response, 1);

  const message = finalizeMessage(response);
  socket.send(message);
}

export function sendChatInput(socket, data) {
  let response = createMessage(oHeadChatInput);
  response = writeString(data.input, response, 1);

  const message = finalizeMessage(response);
  socket.send(message);
}

export function readCompleteUpdate(dataView, offset) {
  const users = {};
  const messages = [];
  let result;

    // Read number of active users
  result = readUint32(dataView, offset);
  offset = result.offset;

    // Read active user ID-username mappings
  const userCount = result.value;
  for( let i = 0; i < userCount; i++ ) {
    result = readUserId(dataView, offset);
    offset = result.offset;

    const userId = result.value;

    result = readString(dataView, offset);
    offset = result.offset;

    const username = result.value;
    users[userId] = username;
  }

    // Read chat messages snapshot
  result = readUint32(dataView, offset);
  const messageCount = result.value;
  offset = result.offset;

  console.log(messageCount)

  for( let i = 0; i < messageCount; i++ ) {
    result = readUserId(dataView, offset);
    const userId = result.value;
    offset = result.offset;

    result = readString(dataView, offset);
    const message = result.value;
    messages.push({userId, message});
    offset = result.offset;
  }

  return {
    users,
    messages
  };
}
