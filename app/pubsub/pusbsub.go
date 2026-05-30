package pubsub

import "net"

type Pubsub struct {
	chanels map[string]*Channel
}

type Channel struct {
	Name        string
	Connections []net.Conn
}

func New() *Pubsub {
	return &Pubsub{
		chanels: make(map[string]*Channel),
	}
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name: name,
	}
}

func (ps *Pubsub) Subscribe(channel *Channel, conn net.Conn) {
	channel.Connections = append(channel.Connections, conn)
	ps.chanels[channel.Name] = channel
}

func (ps *Pubsub) GetChannels() map[string]*Channel {
	return ps.chanels
}

func (ps *Pubsub) GetChannel(name string) *Channel {
	chanel, ok := ps.chanels[name]
	if !ok {
		return nil
	}

	return chanel
}
