package gostun

type transactionIDSetter struct{}

var TransactionID MsgSetter = transactionIDSetter{}

// Sets Message
type MsgSetter interface {
	AddTo(m *Message) error
}
