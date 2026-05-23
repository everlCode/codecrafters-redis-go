package database

type CommandQueue struct {
	Name string
	Args []string
}

func NewCommandQueue(name string, args []string) CommandQueue {
	return CommandQueue{
		Name: name,
		Args: args,
	}
}