package command

const (
	typeFailed = iota
	typeNull
	typeBool
	typeNumber
	typeString
	typeAmbiguous
)

type argument struct {
	argType int
	data    *string
}

func newArgument(argType int, data *string) *argument {
	return &argument{
		argType: argType,
		data:    data,
	}
}

func (arg *argument) ArgType() int {
	return arg.argType
}

func (arg *argument) Data() *string {
	return arg.data
}
