package chat

type Message struct {
	user    *User
	count   int
	message *string
}

func newMessage(user *User, count int, msg *string) *Message {
	return &Message{
		user:    user,
		count:   count,
		message: msg,
	}
}

func (m *Message) User() *User {
	return m.user
}

func (m *Message) Message() *string {
	return m.message
}
