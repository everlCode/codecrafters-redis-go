package clients

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/pubsub"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Client struct {
	connection       net.Conn
	parser           *resp.Parser
	offset           int
	inTx             bool
	txQueue          []CommandQueue
	replica          bool
	RDBSent          bool
	masterConnection bool
	channels    map[string]*pubsub.Channel
}

func New(conn net.Conn) *Client {
	parser := resp.New(conn)
	return &Client{
		connection: conn,
		parser:     parser,
		channels: make(map[string]*pubsub.Channel),
	}
}

func (c *Client) Subscribe(channel *pubsub.Channel) {
	c.channels[channel.Name] = channel
}

func (c *Client) Unbscribe(channelName string) {
	_, ok := c.channels[channelName]
	if !ok {
		return
	}
	delete(c.channels, channelName)
}

func (c *Client) GetSubscribtions() map[string]*pubsub.Channel {
	return c.channels
}

func (c *Client) IsSubscriber() bool {
	return len(c.channels) > 0
}

func (c *Client) SetOffset(v int) {
	c.parser.SetOffset(v)
	c.offset = v
}

func (c *Client) AddOffset(v int) {
	c.offset += v
	c.parser.SetOffset(0)
}

func (c Client) GetOffset() int {
	return c.offset
}

func (c *Client) SetParser(p *resp.Parser) {
	c.parser = p
}

func (c *Client) GetParser() *resp.Parser {
	return c.parser
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
