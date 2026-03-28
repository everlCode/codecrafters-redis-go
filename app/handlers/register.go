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
