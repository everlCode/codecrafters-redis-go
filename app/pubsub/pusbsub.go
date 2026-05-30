package pubsub

import "net"

type Pubsub struct {
	subscribtions []*Subscribtion
}

type Subscribtion struct {
	Conn net.Conn
}

func New() *Pubsub {
	return &Pubsub{}
}

func NewSubscription(conn net.Conn) *Subscribtion {
	return &Subscribtion{
		Conn: conn,
	}
}

func (ps *Pubsub) AddSubscription(subscription *Subscribtion) {
	ps.subscribtions = append(ps.subscribtions, subscription)
}

func (ps *Pubsub) GetSubscription() []*Subscribtion {
	return ps.subscribtions
}