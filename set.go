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

func (t MessageType) SetTo(m *Message) error { return nil }

func (m *Message) NewTransaction() error {
	_, err := io.ReadFunc(rand.Reader, m.TransactionID[:])
	if err != nil {
		return err
	}
	m.WriteTransactionID()
	return nil
}
