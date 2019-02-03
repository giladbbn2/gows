package gows

import "net/http"

type BaseControllerInterface interface {
	SetDirs()

	SetRequest(r *http.Request)
	SetResponse(w http.ResponseWriter)

	GetTemplate(tpl string) ([]byte, error)
}
