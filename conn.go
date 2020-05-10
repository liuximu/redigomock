package redigomock

import "github.com/gomodule/redigo/redis"

var _ redis.Conn = &conn{}

type conn struct {
	mock *redigoMock
}

func (c *conn) do(opt *Operation) (*Value, error) {
	rsp, err := c.mock.Match(opt)
	if err != nil {
		return nil, err
	}

	if err := rsp.do(); err != nil {
		return nil, err
	}

	return rsp.value(), nil
}

// Close closes the connection.
func (c *conn) Close() error {
	val, err := c.do(&Operation{
		Opt: "Close",
	})
	if err != nil {
		return err
	}
	return val.err
}

// Err returns a non-nil value when the connection is not usable.
func (c *conn) Err() error {
	val, err := c.do(&Operation{
		Opt: "Err",
	})
	if err != nil {
		return err
	}
	return val.err
}

// Do sends a command to the server and returns the received reply.
func (c *conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	val, err := c.do(&Operation{
		Opt:  "Do",
		Cmd:  commandName,
		Args: args,
	})
	if err != nil {
		return nil, err
	}
	return val.reply, val.err
}

// Send writes the command to the client's output buffer.
func (c *conn) Send(commandName string, args ...interface{}) error {
	val, err := c.do(&Operation{
		Opt:  "Send",
		Cmd:  commandName,
		Args: args,
	})
	if err != nil {
		return err
	}
	return val.err
}

// Flush flushes the output buffer to the Redis server.
func (c *conn) Flush() error {
	val, err := c.do(&Operation{
		Opt: "Flush",
	})
	if err != nil {
		return err
	}
	return val.err
}

// Receive receives a single reply from the Redis server
func (c *conn) Receive() (reply interface{}, err error) {
	val, err := c.do(&Operation{
		Opt: "Receive",
	})
	if err != nil {
		return nil, err
	}
	return val.reply, val.err
}
