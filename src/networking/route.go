package networking

import (
	"github.com/gorilla/websocket"
	"net/http"
	"tools"
)

//container for api methods (gloabl access for anyone using networking package)
var APILogConnection http.ResponseWriter
var APILogRequest *http.Request
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//create type for consistent function signatures in routing map
type fn func(http.ResponseWriter, *http.Request)

func APIRouteMap() map[string]fn {
	return map[string]fn{
		"/h2":            Info,
		"/api/log":       Logging,
		"/api/logs":      EchoLog,
		"/api/transform": tools.RemodelJ,
	}
}

func UIRouteMap() map[string]string {
	return map[string]string{
		"/sb-admin-btsrp-temp": "views/startbootstrap-sb-admin/index.html",
	}
}
