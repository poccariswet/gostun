package gostun

import (
	"errors"
	"log"
	"sync"
	"time"
)

type CallbackHandle struct {
	Cond     *sync.Cond
	process  bool
	callback func(e EventObject)
}

var callbackPool = sync.Pool{
	New: func() interface{} {
		return &CallbackHandle{
			Cond: sync.NewCond(new(sync.Mutex)),
		}
	},
}

func (c *CallbackHandle) Reset() {
	c.process = false
	c.callback = nil
}

func (c *CallbackHandle) HandleEvent(e EventObject) {
	if c.callback == nil {
		log.Fatal("callback is nil")
	}
	c.callback(e)
	c.Cond.L.Lock()
	c.process = true
	c.Cond.Broadcast()
	c.Cond.L.Unlock()
}

func (c *CallbackHandle) Wait() {
	c.Cond.L.Lock()
	for !c.process {
		c.Cond.Wait()
	}
	c.Cond.L.Unlock()
}

func (a *Agent) TransactionHandle(id [TransactionIDSize]byte, h Handler, rto time.Time) error {
	a.mux.Lock()
	defer a.mux.Unlock()

	if a.closed {
		return errors.New("agent closed")
	}

	_, exist := a.transactions[id]
	if exist {
		return errors.New("transaction exists with same id")
	}

	a.transactions[id] = TransactionAgent{
		id:      id,
		handler: h,
		Timeout: rto,
	}

	return nil
}

func (c *Client) TransactionLaunch(m *Message, h Handler, rto time.Time) error {
	c.rw.RLock()
	closed := c.clientclose
	c.rw.RUnlock()

	if closed {
		return errors.New("client closed")
	}

	if h != nil {
		if err := c.agent.TransactionHandle(m.TransactionID, h, rto); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Call(m *Message, h func(EventObject), rto time.Time) error {
	f := callbackPool.Get().(*CallbackHandle)
	f.callback = h

	defer func() {
		f.Reset()
		callbackPool.Put(f)
	}()

	// waiting TransactionLaunch until call callback func
	if err := c.TransactionLaunch(m, f, rto); err != nil {
		return err
	}
	f.Wait()

	return nil
}
