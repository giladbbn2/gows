package gows

import (
	"net/http"
)

type BaseControllerInterface interface {
	SetWebServer(ws *WebServer)
	GetWebServer() *WebServer

	GetMicroServiceRootDir() string

	SetRequest(r *http.Request)
	SetResponse(w http.ResponseWriter)
	GetResponse() *http.ResponseWriter
}
