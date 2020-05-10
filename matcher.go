package redigomock

import (
	"fmt"
	"reflect"
	"strings"
)

type Operation struct {
	Opt  string
	Cmd  string
	Args []interface{}
}

func (opt *Operation) Equal(b *Operation, config *Config) bool {
	if b == nil {
		return false
	}

	if opt.Opt != b.Opt {
		return false
	}

	if strings.ToLower(opt.Cmd) != strings.ToLower(b.Cmd) {
		return false
	}

	if config == nil || config.FuzzyMatch {
		return true
	}

	return equal(opt.Args, b.Args)
}

func (opt *Operation) String() string {
	return fmt.Sprintf("Opt[%s] Cmd[%s] Args[%v]", opt.Opt, opt.Cmd, opt.Args)
}

func equal(as, bs []interface{}) bool {
	if len(as) == 0 && len(bs) == 0 {
		return true
	}

	if len(as) != len(bs) {
		return false
	}

	for i, a := range as {
		b := bs[i]
		if !reflect.DeepEqual(a, b) {
			return false
		}
	}

	return true
}

type Expection struct {
	triggered bool

	operation *Operation

	reply Action
}

type Matcher interface {
	Match(operation *Operation) (Action, error)
}

var _ Matcher = &redigoMock{}

func (rm *redigoMock) Match(operation *Operation) (Action, error) {
	rm.Lock()
	defer rm.Unlock()

	for _, expection := range rm.expections {
		if expection.triggered {
			continue
		}

		if expection.operation.Equal(operation, rm.config) {
			expection.triggered = true
			return expection.reply, nil
		}

		if rm.config.Order {
			rm.matchErr = fmt.Errorf("want: %s, got: %s", expection.operation, operation)
			return nil, rm.matchErr
		}
	}

	rm.matchErr = fmt.Errorf("got: %s, cant find operation", operation)
	return nil, rm.matchErr
}
