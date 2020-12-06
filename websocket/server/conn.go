package server

import (
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"sync"
	"time"
)

type Conn struct {
	Conn            *websocket.Conn
	AfterReadFunc   func(messageType int, r io.Reader)
	BeforeCloseFunc func()
	once            sync.Once
	id              string
	stopCh          chan struct{}
}

func (c *Conn) Write(p []byte) (n int, err error) {
	select {
	case <-c.stopCh:
		return 0, errors.New("conn is closed")
	default:
		err = c.Conn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			return 0, err
		}
		return len(p), nil
	}
}

func (c *Conn) GetID() string {
	c.once.Do(func() {
		u := "ddd"
		c.id = u
	})
	return c.id
}

func (c *Conn) Close() error {
	select {
	case <-c.stopCh:
		return errors.New("Conn already been closed")
	default:
		c.Conn.Close()
		close(c.stopCh)
		return nil
	}
}

func (c *Conn) Listen() {
	c.Conn.SetCloseHandler(func(code int, text string) error {
		if c.BeforeCloseFunc != nil {
			c.BeforeCloseFunc()
		}
		if err := c.Close(); err != nil {
			log.Println(err)
		}
		message := websocket.FormatCloseMessage(code, "")
		_ = c.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})
ReadLoap:
	for {
		select {
		case <-c.stopCh:
			break ReadLoap
		default:
			messageType, r, err := c.Conn.NextReader()
			if err != nil {
				break ReadLoap
			}
			if c.AfterReadFunc != nil {
				c.AfterReadFunc(messageType, r)
			}
		}
	}
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{
		Conn:   conn,
		stopCh: make(chan struct{}),
	}
}
