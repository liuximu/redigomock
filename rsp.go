package redigomock

import (
	"time"
)

type basicRspStub interface {
	WillPanic(pani interface{})
	WillDelay(dur time.Duration)
}

type ErrRspStub interface {
	basicRspStub
	WillReturnError(err error)
}

type RelayRspStub interface {
	ErrRspStub
	WillRelay(relay interface{})
}

type Action interface {
	do() error
	value() *Value
}

type Value struct {
	err   error
	relay interface{}
}

var (
	_ Action = &basicRsp{}
	_ Action = &errorRsp{}
	_ Action = &relayRsp{}

	_ basicRspStub = &basicRsp{}
	_ ErrRspStub   = &errorRsp{}
	_ RelayRspStub = &relayRsp{}
)

type basicRsp struct {
	delay time.Duration
	pani  interface{}
}

func (dr *basicRsp) WillPanic(pani interface{}) {
	dr.pani = pani
}

func (dr *basicRsp) WillDelay(dur time.Duration) {
	dr.delay = dur
}

func (dr *basicRsp) value() *Value {
	panic("never invoke me")
}

func (dr *basicRsp) do() error {
	if dr.pani != nil {
		panic(dr.pani)
	}

	if dr.delay != 0 {
		time.Sleep(dr.delay)
	}
	return nil
}

type errorRsp struct {
	basicRsp
	err error
}

func (er *errorRsp) WillReturnError(err error) {
	er.err = err
}
func (er *errorRsp) value() *Value {
	return &Value{
		err: er.err,
	}
}

func (er *errorRsp) do() error {
	return nil
}

type relayRsp struct {
	errorRsp
	relay interface{}
}

func (rr *relayRsp) WillRelay(relay interface{}) {
	rr.relay = relay
}

func (rr *relayRsp) value() *Value {
	return &Value{
		err:   rr.err,
		relay: rr.relay,
	}
}
