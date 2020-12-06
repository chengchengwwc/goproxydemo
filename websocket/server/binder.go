package server

import (
	"errors"
	"fmt"
	"sync"
)

type eventConn struct {
	Event string
	Conn  *Conn
}

type binder struct {
	mu                  sync.RWMutex
	userID2EventConnMap map[string]*[]eventConn
	connID2UserIDMap    map[string]string
}

func (b *binder) Bind(userID, event string, conn *Conn) error {
	if userID == "" {
		return errors.New("userID can't be empty")
	}
	if event == "" {
		return errors.New("event can't be empty")
	}

	if conn == nil {
		return errors.New("bad")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if eConns, ok := b.userID2EventConnMap[userID]; ok {
		for i := range *eConns {
			if (*eConns)[i].Conn == conn {
				return nil
			}
		}
		newEConns := append(*eConns, eventConn{event, conn})
		b.userID2EventConnMap[userID] = &newEConns
	} else {
		b.userID2EventConnMap[userID] = &[]eventConn{{event, conn}}
	}
	b.connID2UserIDMap[conn.GetID()] = userID
	return nil
}

func (b *binder) Unbind(conn *Conn) error {
	if conn == nil {
		return errors.New("bad")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	userID, ok := b.connID2UserIDMap[conn.GetID()]
	if !ok {
		return errors.New("xxxx")
	}
	if eConns, ok := b.userID2EventConnMap[userID]; ok {
		for i := range *eConns {
			if (*eConns)[i].Conn == conn {
				newEConns := append((*eConns)[:i], (*eConns)[i+1:]...)
				b.userID2EventConnMap[userID] = &newEConns
				delete(b.connID2UserIDMap, conn.GetID())

				// delete the key of userID when the length of the related
				// eventConn slice is 0.
				if len(newEConns) == 0 {
					delete(b.userID2EventConnMap, userID)
				}
				return nil
			}
		}
		return fmt.Errorf("can't find the conn of ID: %s", conn.GetID())
	}
	return fmt.Errorf("can't find the eventConns by userID: %s", userID)
}

func (b *binder) FindConn(connID string) (*Conn, bool) {
	if connID == "" {
		return nil, false
	}
	userID, ok := b.connID2UserIDMap[connID]
	if ok {
		if eConns, ok := b.userID2EventConnMap[userID]; ok {
			for i := range *eConns {
				if (*eConns)[i].Conn.GetID() == connID {
					return (*eConns)[i].Conn, true
				}
			}
		}
		return nil, false
	}
	for _, eConns := range b.userID2EventConnMap {
		for i := range *eConns {
			if (*eConns)[i].Conn.GetID() == connID {
				return (*eConns)[i].Conn, true
			}
		}
	}
	return nil, true
}

func (b *binder) FilterConn(userID, event string) ([]*Conn, error) {
	if userID == "" {
		return nil, errors.New("bad")
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	if eConns, ok := b.userID2EventConnMap[userID]; ok {
		esc := make([]*Conn, 0, len(*eConns))
		for i := range *eConns {
			if event == "" || (*eConns)[i].Event == event {
				esc = append(esc, (*eConns)[i].Conn)
			}
		}
		return esc, nil
	}
	return []*Conn{}, nil

}
