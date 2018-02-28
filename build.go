package gostun

import "log"

type transactionIDSetter struct{}

func (transactionIDSetter) AddTo(m *Message) error {
	return nil
}

// Sets *Message attribute.
type MsgSetter interface {
	AddTo(m *Message) error
}

var TransactionID MsgSetter = transactionIDSetter{}

// reset message
func (m *Message) Reset() {
	m.Raw = m.Raw[:0]
	m.Length = 0
	m.Attributes = m.Attributes[:0]
}

func (m *Message) AllocRaw() {
	l := len(m.Raw) + messageHeader
	for cap(m.Raw) < l {
		m.Raw = append(m.Raw, 0)
	}
	m.Raw = m.Raw[:l]
}

// make message header
func (m *Message) WriteMessageHeader() {
	m.AllocRaw() // alloc 0, part of message header size

}

func (m *Message) build(s ...MsgSetter) error {
	m.Reset()
	m.WriteMessageHeader()

	return nil
}

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
