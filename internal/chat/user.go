package chat

type UserID uint64

type User struct {
	id   UserID
	name *string
}

func NewUser(id UserID, name *string) *User {
	return &User{
		id:   id,
		name: name,
	}
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Name() *string {
	return u.name
}

func (u *User) SetName(name *string) {
	u.name = name
}
