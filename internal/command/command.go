package command

type command struct {
	name      string
	arguments []*argument
}

func (comm *command) Name() *string {
	return &comm.name
}

func (comm *command) GetArgument(index int) *argument {
	return comm.arguments[index]
}
