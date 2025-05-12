package command

type Command struct {
	name      string
	arguments []*Argument
}

func (comm *Command) Name() *string {
	return &comm.name
}

func (comm *Command) GetArgument(index int) *Argument {
	return comm.arguments[index]
}
