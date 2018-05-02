package gostun

import (
	"crypto/rand"
	"io"
)

type SetTransaer struct{}

var TransactionID Transaer = SetTransaer{}

// Sets Message attr
type Transaer interface {
	SetTo(m *Message) error
}

func (SetTransaer) SetTo(m *Message) error {
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
