package server

import (
	"bytes"
	"encoding/binary"
)

type messageHeader uint8
type userID uint64

const (
	// Inbound message header (from client)
	iHeadClientInfo = iota + 1 // Client is sending its info upon joining.
	iHeadNameChange            // Client is requesting a name change.

	// Outbound message header (to client)
	oHeadActiveUsers   // Server is sending a list of active users.
	oHeadInactiveUsers // Server is sending a list of users no longer active.
	oHeadSnapshot      // Server is sending the latest snapshot of chat messages (can be complete or delta).
)

const (
	firstUserID userID = 1
)

type message struct {
	header messageHeader
	body   []byte
}

func writeUserId(id userID, buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, id)
}

func writeString(str string, buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, uint32(len(str)))
	buffer.WriteString(str)
}

func writeChatMessage(id userID, chat string, buffer *bytes.Buffer) {
	writeUserId(id, buffer)
	writeString(chat, buffer)
}

func newActiveUsers(id []userID, name []string) *message {
	buffer := new(bytes.Buffer)

	for i := range len(id) {
		writeUserId(id[i], buffer)
		writeString(name[i], buffer)
	}

	return &message{
		header: oHeadActiveUsers,
		body:   buffer.Bytes(),
	}
}

func newInactiveUsers(id []userID) *message {
	buffer := new(bytes.Buffer)

	for _, userID := range id {
		writeUserId(userID, buffer)
	}

	return &message{
		header: oHeadInactiveUsers,
		body:   buffer.Bytes(),
	}
}
