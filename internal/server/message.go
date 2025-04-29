package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// import (
// 	"bytes"
// 	"encoding/binary"
// )

// type netHeader uint8
// type userID uint64

// const (
// 	firstUserID userID = 1
// )

// type netMessage struct {
// 	head netHeader
// 	body []byte
// }

// type user struct {
// 	id   userID
// 	name string
// }

// type chatMessage struct {
// 	user    user
// 	message string
// }

// type snapshot struct {
// 	messages []chatMessage
// }

// func newUser(id userID, name *string) *user {
// 	return &user{
// 		id:   id,
// 		name: *name,
// 	}
// }

// func writeUserId(id userID, buffer *bytes.Buffer) {
// 	binary.Write(buffer, binary.LittleEndian, id)
// }

// func writeString(str string, buffer *bytes.Buffer) {
// 	binary.Write(buffer, binary.LittleEndian, uint32(len(str)))
// 	buffer.WriteString(str)
// }

// func writeChatMessage(message *chatMessage, buffer *bytes.Buffer) {
// 	writeUserId(message.user.id, buffer)
// 	writeString(message.message, buffer)
// }

// func newActiveUsers(users []user) *netMessage {
// 	buffer := new(bytes.Buffer)

// 	for _, user := range users {
// 		writeUserId(user.id, buffer)
// 		writeString(user.name, buffer)
// 	}

// 	return &netMessage{
// 		head: oHeadActiveUsers,
// 		body: buffer.Bytes(),
// 	}
// }

// func newInactiveUsers(ids []userID) *netMessage {
// 	buffer := new(bytes.Buffer)

// 	for _, userID := range ids {
// 		writeUserId(userID, buffer)
// 	}

// 	return &netMessage{
// 		head: oHeadInactiveUsers,
// 		body: buffer.Bytes(),
// 	}
// }

// func newSnapshot(messages []chatMessage) *netMessage {
// 	buffer := new(bytes.Buffer)

// 	for _, message := range messages {
// 		writeChatMessage(&message, buffer)
// 	}

// 	return &netMessage{
// 		head: oHeadSnapshot,
// 		body: buffer.Bytes(),
// 	}
// }

type header uint8

type request struct {
	head   header
	client *client
	bytes  []byte
}

const firstUserID userID = 1
const stringDelim byte = 0

const (
	// Inbound message header (from client)
	iHeadClientInfo = iota + 1 // Client is sending its info upon joining.
	iHeadNameChange            // Client is requesting a name change.
	iHeadChatInput             // Client is sending a chat message/command

	// Outbound message header (to client)
	oHeadActiveUsers   // Server is sending a list of active users.
	oHeadInactiveUsers // Server is sending a list of users no longer active.
	oHeadSnapshot      // Server is sending the latest snapshot of chat messages (can be complete or delta).
)

func newUserInfo(id userID, name *string) *userInfo {
	return &userInfo{
		id:   id,
		name: name,
	}
}

func readString(buffer *bytes.Buffer) *string {
	line, err := buffer.ReadString(stringDelim)

	if err != nil {
		fmt.Println("ERROR: Unable to read string from a client message!")
		fmt.Println(err.Error())
	}

	return &line
}

func writeString(str string, buffer *bytes.Buffer) {
	buffer.WriteString(str)
	buffer.WriteByte(stringDelim)
}

func writeUserId(id userID, buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, id)
}

func createResponse(head header) *bytes.Buffer {
	res := new(bytes.Buffer)
	res.WriteByte((uint8)(head))
	return res
}
