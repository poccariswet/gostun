package gostun

import "log"

type transactionIDSetter struct{}

// Sets *Message attribute.
type MsgSetter interface {
	AddTo(m *Message) error
}

var TransactionID MsgSetter = transactionIDSetter{}

func MessageBuild(s ...MsgSetter) *Message {
	m, err := Build(s...)
	if err != nil {
		log.Fatal(err)
	}
}
