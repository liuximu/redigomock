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

type ReplyRspStub interface {
	ErrRspStub
	WillReply(reply interface{})
}

type Action interface {
	do() error
	value() *Value
}

type Value struct {
	err   error
	reply interface{}
}

var (
	_ Action = &basicRsp{}
	_ Action = &errorRsp{}
	_ Action = &replyRsp{}

	_ basicRspStub = &basicRsp{}
	_ ErrRspStub   = &errorRsp{}
	_ ReplyRspStub = &replyRsp{}
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

type replyRsp struct {
	errorRsp
	reply interface{}
}

func (rr *replyRsp) WillReply(reply interface{}) {
	rr.reply = reply
}

func (rr *replyRsp) value() *Value {
	return &Value{
		err:   rr.err,
		reply: rr.reply,
	}
}
