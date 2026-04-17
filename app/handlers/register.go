package handlers

import (
	"errors"
)

type Register struct {
	commands map[string]Command
}

func NewRegister() *Register {
	register := &Register{
		commands: make(map[string]Command),
	}
	register.Add(&PingCommand{})
	register.Add(&SetCommand{})
	register.Add(&GetCommand{})
	register.Add(&EchoCommand{})
	register.Add(&LpushCommand{})
	register.Add(&RpushCommand{})
	register.Add(&LRangeCommand{})
	register.Add(&LLenCommand{})
	register.Add(&LPopCommand{})
	register.Add(&BlPopCommand{})

	return register
}

func (r *Register) Add(command Command) {
	r.commands[command.Name()] = command
}

func (r *Register) Get(name string) (Command, error) {
	command, ok := r.commands[name]
	if !ok {
		return nil, errors.New("Undefined command!")
	}

	return command, nil
}
