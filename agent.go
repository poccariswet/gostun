package gostun

import (
	"errors"
	"sync"
	"time"
)

// process of transaction in message
type Agent struct {
	transactions map[transactionID]TransactionAgent
	mux          sync.Mutex
	nonHandler   Handler // non-registered transactions
	closed       bool
}

type transactionID [TransactionIDSize]byte //12byte, 96bit

// transaction in progress
type TransactionAgent struct {
	ID      transactionID
	Timeout time.Time
	handler Handler // if transaction is succeed will be called
}

type AgentHandle struct {
	handler Handler
}

// reference http.HandlerFunc same work
type Handler interface {
	HandleEvent(e MessageObj)
}

// type HandleFunc func(e EventObject)

// func (f HandleFunc) HandleEvent(e EventObject) {
// 	f(e)
// }

type MessageObj struct {
	Msg *Message
	Err error
}

func NewAgent() *Agent {
	h := AgentHandle{}
	a := &Agent{
		transactions: make(map[transactionID]TransactionAgent),
		nonHandler:   h.handler,
	}
	return a
}

func (a *Agent) ProcessHandle(m *Message) error {
	e := MessageObj{
		Msg: m,
	}

	a.mux.Lock() // protect to do multiple access to transaction
	tr, ok := a.transactions[m.TransactionID]
	delete(a.transactions, m.TransactionID) //delete maps entry
	a.mux.Unlock()

	if ok {
		tr.handler.HandleEvent(e) // HandleEvent cast the e to hander type
	} else if a.nonHandler != nil {
		a.nonHandler.HandleEvent(e) // the transaction is not registered
	}
	return nil
}

/*
すべてのハンドラがTransactionTimeOutErrを処理するまで、
指定された時刻より前にデッドラインを持つすべてのトランザクションをblockする。
エージェントが既に閉じられている場合、ErrAgentを返す
*/

var (
	ErrAgent              = errors.New("agent closed")
	TransactionTimeOutErr = errors.New("transaction is timed out")
)

/*
The value for RTO SHOULD be cached by a client after the completion
of the transaction, and used as the starting value for RTO for the
next transaction to the same server (based on equality of IP
address).
*/

func (a *Agent) TimeOutHandle(trate time.Time) error {
	call := make([]Handler, 0, 100)
	remove := make([]transactionID, 0, 100)
	a.mux.Lock()

	if a.closed {
		a.mux.Unlock()
		return ErrAgent
	}

	for i, tr := range a.transactions {
		if tr.Timeout.Before(trate) {
			call = append(call, tr.handler)
			remove = append(remove, i)
		}
	}

	// no registered transactions
	for _, id := range remove {
		delete(a.transactions, id)
	}

	a.mux.Unlock()
	e := MessageObj{
		Err: TransactionTimeOutErr,
	}
	// return transactions
	for _, h := range call {
		h.HandleEvent(e)
	}

	return nil
}
