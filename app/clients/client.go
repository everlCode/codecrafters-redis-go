package clients

import "net"

type Client struct {
	connection net.Conn
	inTx    bool
	txQueue []CommandQueue
	replica bool
	RDBSent bool
	masterConnection bool
}

func New(conn net.Conn) *Client {
	return &Client{
		connection: conn,
	}
}

func (c *Client) GetConnection() net.Conn {
	return c.connection
}

func (c *Client) IsTransaction() bool {
	return c.inTx
}

func (c *Client) StartTransactions() {
	c.inTx = true
}

func (c *Client) EndTransactions() {
	c.inTx = false
}

func (c *Client) GetCommandQueue() []CommandQueue {
	return c.txQueue
}

func (c *Client) PushCommandQueue(v CommandQueue) {
	c.txQueue = append(c.txQueue, v)
}

func (c *Client) ClearCommandQueue() {
	c.txQueue = []CommandQueue{}
}

func (c *Client) IsReplica() bool {
	return c.replica
}

func (c *Client) SetReplica(v bool) {
	c.replica = v
}

func (c *Client) SetMasterConnection(v bool) {
	c.masterConnection = v
}

func (c *Client) MasterConnection() bool {
	return c.masterConnection
}
