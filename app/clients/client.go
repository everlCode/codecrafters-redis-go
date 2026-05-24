package clients

type Client struct {
	inTx    bool
	txQueue []CommandQueue
}

func New() *Client {
	return &Client{}
}

func (c Client) IsTransaction() bool {
	return c.inTx
}

func (c *Client) StartTransactions() {
	c.inTx = true
}

func (c *Client) EndTransactions() {
	c.inTx = false
}

func (c Client) GetCommandQueue() []CommandQueue {
	return c.txQueue
}

func (c *Client) PushCommandQueue(v CommandQueue) {
	c.txQueue = append(c.txQueue, v)
}

func (c Client) ClearCommandQueue() {
	c.txQueue = []CommandQueue{}
}
