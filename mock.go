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
	ExceptClose() ErrRspStub // error

	// Err returns a non-nil value when the connection is not usable.
	ExceptInvokeErr() ErrRspStub // error

	// Do sends a command to the server and returns the received reply.
	ExceptDo(commandName string, args ...interface{}) RelayRspStub // (reply interface{}, err error)

	// Send writes the command to the client's output buffer.
	ExceptSend(commandName string, args ...interface{}) ErrRspStub // error

	// Flush flushes the output buffer to the Redis server.
	ExceptFlush() ErrRspStub

	// Receive receives a single reply from the Redis server
	ExceptReceive() RelayRspStub // (reply interface{}, err error)
}

type redigoMock struct {
	sync.Mutex

	matchErr error
	config   *Config

	exceptions []*Exception
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

	for _, exception := range rm.exceptions {
		if !exception.triggered {
			return fmt.Errorf("%s not be triggered", exception.operation)
		}
	}

	return nil
}

// Close closes the connection.
func (rm *redigoMock) ExceptClose() ErrRspStub {
	rsp := &errorRsp{}
	rm.exceptions = append(rm.exceptions, &Exception{
		operation: &Operation{
			Opt: "Close",
		},
		relay: rsp,
	})

	return rsp
}

// Err returns a non-nil value when the connection is not usable.
func (rm *redigoMock) ExceptInvokeErr() ErrRspStub {
	rsp := &errorRsp{}
	rm.exceptions = append(rm.exceptions, &Exception{
		operation: &Operation{
			Opt: "Err",
		},
		relay: rsp,
	})

	return rsp
}

// Do sends a command to the server and returns the received reply.
func (rm *redigoMock) ExceptDo(commandName string, args ...interface{}) RelayRspStub {
	rsp := &relayRsp{}
	rm.exceptions = append(rm.exceptions, &Exception{
		operation: &Operation{
			Opt:  "Do",
			Cmd:  commandName,
			Args: args,
		},
		relay: rsp,
	})

	return rsp
}

// Send writes the command to the client's output buffer.
func (rm *redigoMock) ExceptSend(commandName string, args ...interface{}) ErrRspStub {
	rsp := &errorRsp{}
	rm.exceptions = append(rm.exceptions, &Exception{
		operation: &Operation{
			Opt:  "Send",
			Cmd:  commandName,
			Args: args,
		},
		relay: rsp,
	})

	return rsp
}

// Flush flushes the output buffer to the Redis server.
func (rm *redigoMock) ExceptFlush() ErrRspStub {
	rsp := &errorRsp{}
	rm.exceptions = append(rm.exceptions, &Exception{
		operation: &Operation{
			Opt: "Flush",
		},
		relay: rsp,
	})

	return rsp
}

// Receive receives a single reply from the Redis server
func (rm *redigoMock) ExceptReceive() RelayRspStub {
	rsp := &relayRsp{}
	rm.exceptions = append(rm.exceptions, &Exception{
		operation: &Operation{
			Opt: "Receive",
		},
		relay: rsp,
	})

	return rsp
}
