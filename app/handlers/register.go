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
	register.Add(PING, &PingCommand{})
	register.Add(SET, &SetCommand{})
	register.Add(GET, &GetCommand{})
	register.Add(ECHO, &EchoCommand{})
	register.Add(LPUSH, &LpushCommand{})
	register.Add(RPUSH, &RpushCommand{})
	register.Add(LRANGE, &LRangeCommand{})
	register.Add(LLEN, &LLenCommand{})
	register.Add(LPOP, &LPopCommand{})
	register.Add(BLPOP, &BlPopCommand{})
	register.Add(TYPE, &TypeCommand{})
	register.Add(XADD, &XaddCommand{})
	register.Add(XRANGE, &XrangeCommand{})
	register.Add(XREAD, &XreadCommand{})

	return register
}

func (r *Register) Add(name string, command Command) {
	r.commands[name] = command
}

func (r *Register) Get(name string) (Command, error) {
	command, ok := r.commands[name]
	if !ok {
		return nil, errors.New("Undefined command!")
	}

	return command, nil
}
