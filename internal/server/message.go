package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type header uint8

type request struct {
	head   header
	client *client
	bytes  []byte
}

const firstUserID userID = 1
const stringDelim byte = 0
const empty byte = 0

const (
	// Inbound message header (from client)
	iHeadClientInfo = iota + 1 // Client is sending its info upon joining.
	iHeadNameChange            // Client is requesting a name change.
	iHeadChatInput             // Client is sending a chat message/command

	// Outbound message header (to client)
	oHeadNameChange     // Server is confirming a name change.
	oHeadCompleteUpdate // Server is sending a complete update containing latest messages, active and inactive users.
	oHeadDeltaUpdate    // Server is sending the latest delta snapshot of messages, active and inactive users.
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

func writeEmpty(buffer *bytes.Buffer) {
	buffer.WriteByte(empty)
}

func writeString(str string, buffer *bytes.Buffer) {
	buffer.WriteString(str)
	buffer.WriteByte(stringDelim)
}

func writeUserId(id userID, buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, id)
}

func writeUInt32(n uint32, buffer *bytes.Buffer) {
	binary.Write(buffer, binary.LittleEndian, n)
}

func writeUserInfo(u *userInfo, b *bytes.Buffer) {
	if u != nil {
		writeUserId(u.id, b)
		writeString(*u.name, b)
	} else {
		writeUserId(0, b)
	}
}

func createResponse(head header) *bytes.Buffer {
	res := new(bytes.Buffer)
	res.WriteByte((uint8)(head))
	return res
}
