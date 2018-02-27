package gostun

import "log"

type transactionIDSetter struct{}

// Sets *Message attribute.
type MsgSetter interface {
	AddTo(m *Message) error
}

var TransactionID MsgSetter = transactionIDSetter{}

func (m *Message) build(s ...MsgSetter) error {}

// wraps m.build
func Build(s ...MsgSetter) (*Message, error) {
	m := new(Message)
	return m, m.build(s...)
}

func MessageBuild(s ...MsgSetter) *Message {
	m, err := Build(s...)
	if err != nil {
		log.Fatal(err)
	}

	return m
}
