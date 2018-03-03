package main

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
	for !c.processed {
		c.cond.Wait()
	}
	c.cond.L.Unlock()
}

func (c *Client) TransactionLaunch(m *Message, h Handler, rto time.Time) error {
	c.rw.RLock()
	closed := c.close
	c.rw.RUnlock()

	if closed {
		return errors.New("client closed")
	}

	if h != nil {
		if err := c.agent.ProccessLaunch(m.TransactionID, h, rto); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Call(m *Message, h Handler, rto time.Time) error {
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
