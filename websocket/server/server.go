package server

import (
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	serverDefaultWSPath   = "/ws"
	serverDefaultPushPath = "/push"
)

var defaultUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

type Server struct {
	Addr      string
	WSPath    string
	PushPath  string
	Upgrader  *websocket.Upgrader
	AuthToken func(token string) (userID string, ok bool)
	PushAuth  func(r *http.Request) bool
}
