package chat

type ID uint64

type User struct {
	id   ID
	name *string
}

func NewUser(id ID, name *string) *User {
	return &User{
		id:   id,
		name: name,
	}
}

func (u *User) ID() ID {
	return u.id
}

func (u *User) Name() *string {
	return u.name
}

func (u *User) SetName(name *string) {
	u.name = name
}
