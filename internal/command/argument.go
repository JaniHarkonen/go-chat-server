package command

import "strconv"

const (
	TypeFailed = iota
	TypeNull
	TypeBool
	TypeNumber
	TypeString
	TypeUser
	TypeAmbiguous
)

type Argument struct {
	argType int
	data    *string
}

func newArgument(argType int, data *string) *Argument {
	return &Argument{
		argType: argType,
		data:    data,
	}
}

func (arg *Argument) ArgType() int {
	return arg.argType
}

func (arg *Argument) Data() *string {
	return arg.data
}

func (arg *Argument) AsString() string {
	return *arg.data
}

func (arg *Argument) AsInt() (int64, error) {
	return strconv.ParseInt(*arg.data, 10, 0)
}
