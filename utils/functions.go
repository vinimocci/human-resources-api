package utils

import (
	"net/http"
	"github.com/gorilla/websocket"
)

var NetUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}