package redigomock

import (
	"fmt"
	"sync"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Order      bool
	FuzzyMatch bool
}

type Option func(*Config)

// New return one conn and mock
func New(opts ...Option) (redis.Conn, RedigoMock) {
	config := &Config{
		Order:      true,
		FuzzyMatch: true,
	}
	for _, opt := range opts {
		opt(config)
	}
	mock := &redigoMock{
		config: config,
	}
	conn := &conn{
		mock: mock,
	}

	return conn, mock
}

type RedigoMock interface {
	// MatchExpectationsInOrder(bool)

	ExpectationsWereMet() error

	// Close closes the connection.
	ExpectClose() ErrRspStub // error

	// Err returns a non-nil value when the connection is not usable.
	ExpectInvokeErr() ErrRspStub // error

	// Do sends a command to the server and returns the received reply.
	ExpectDo(commandName string, args ...interface{}) ReplyRspStub // (reply interface{}, err error)

	// Send writes the command to the client's output buffer.
	ExpectSend(commandName string, args ...interface{}) ErrRspStub // error

	// Flush flushes the output buffer to the Redis server.
	ExpectFlush() ErrRspStub

	// Receive receives a single reply from the Redis server
	ExpectReceive() ReplyRspStub // (reply interface{}, err error)
}

type redigoMock struct {
	sync.Mutex

	matchErr error
	config   *Config

	expections []*Expection
}

// func (rm *redigoMock) MatchExpectationsInOrder(order bool) {
// 	rm.config.order = order
// }
func (rm *redigoMock) ExpectationsWereMet() error {
	rm.Lock()
	defer rm.Unlock()

	if rm.matchErr != nil {
		return rm.matchErr
	}

	for _, expection := range rm.expections {
		if !expection.triggered {
			return fmt.Errorf("%s not be triggered", expection.operation)
		}
	}

	return nil
}

// Close closes the connection.
func (rm *redigoMock) ExpectClose() ErrRspStub {
	rsp := &errorRsp{}
	rm.expections = append(rm.expections, &Expection{
		operation: &Operation{
			Opt: "Close",
		},
		reply: rsp,
	})

	return rsp
}

// Err returns a non-nil value when the connection is not usable.
func (rm *redigoMock) ExpectInvokeErr() ErrRspStub {
	rsp := &errorRsp{}
	rm.expections = append(rm.expections, &Expection{
		operation: &Operation{
			Opt: "Err",
		},
		reply: rsp,
	})

	return rsp
}

// Do sends a command to the server and returns the received reply.
func (rm *redigoMock) ExpectDo(commandName string, args ...interface{}) ReplyRspStub {
	rsp := &replyRsp{}
	rm.expections = append(rm.expections, &Expection{
		operation: &Operation{
			Opt:  "Do",
			Cmd:  commandName,
			Args: args,
		},
		reply: rsp,
	})

	return rsp
}

// Send writes the command to the client's output buffer.
func (rm *redigoMock) ExpectSend(commandName string, args ...interface{}) ErrRspStub {
	rsp := &errorRsp{}
	rm.expections = append(rm.expections, &Expection{
		operation: &Operation{
			Opt:  "Send",
			Cmd:  commandName,
			Args: args,
		},
		reply: rsp,
	})

	return rsp
}

// Flush flushes the output buffer to the Redis server.
func (rm *redigoMock) ExpectFlush() ErrRspStub {
	rsp := &errorRsp{}
	rm.expections = append(rm.expections, &Expection{
		operation: &Operation{
			Opt: "Flush",
		},
		reply: rsp,
	})

	return rsp
}

// Receive receives a single reply from the Redis server
func (rm *redigoMock) ExpectReceive() ReplyRspStub {
	rsp := &replyRsp{}
	rm.expections = append(rm.expections, &Expection{
		operation: &Operation{
			Opt: "Receive",
		},
		reply: rsp,
	})

	return rsp
}
