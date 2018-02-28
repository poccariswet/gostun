package gostun

import (
	"crypto/rand"
	"io"
)

type transactionIDSetter struct{}

var TransactionID MsgSetter = transactionIDSetter{}

// Sets Message attr
type MsgSetter interface {
	SetTo(m *Message) error
}

func (transactionIDSetter) SetTo(m *Message) error {
	return m.NewTransaction()
}

func (t MessageType) SetTo(m *Message) error {
	m.TypeSet(t)
	return nil
}

func (m *Message) NewTransaction() error {
	_, err := io.ReadFull(rand.Reader, m.TransactionID[:])
	if err != nil {
		return err
	}
	m.WriteTransactionID()
	return nil
}

func (m *Message) TypeSet(t MessageType) {
	m.Type = t           // Message type set
	m.WriteMessageType() // Write Class and Method
}
