package gostun

type transactionIDSetter struct{}

var TransactionID MsgSetter = transactionIDSetter{}

// Sets Message
type MsgSetter interface {
	SetTo(m *Message) error
}

func (transactionIDSetter) SetTo(m *Message) error {
	return m.NewTransactionID()
}
