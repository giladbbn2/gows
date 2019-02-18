package gows

import (
	"errors"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type WebServer struct {
	mux          *http.ServeMux
	addrWithPort string
}

func NewWebServer(addrWithPort string) (*WebServer, error) {

	if addrWithPort == "" {
		return nil, errors.New("parameter can't be zero-valued")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

		io.WriteString(w, "pong")

	})

	ws := new(WebServer)
	ws.mux = mux
	ws.addrWithPort = addrWithPort

	return ws, nil

}

func (ws *WebServer) ListenAndServe() error {

	fs := http.FileServer(http.Dir(ws.GetMicroServiceRootDir() + string(os.PathSeparator) + "includes"))
	http.Handle("/includes", fs)

	err := http.ListenAndServe(ws.addrWithPort, ws.mux)
	if err != nil {
		return err
	}

	return nil

}

func (ws *WebServer) RegisterController(ctrlPattern string, ctrlVer string, ctrl BaseControllerInterface) error {

	if ctrlPattern == "" || ctrlVer == "" || ctrl == nil {
		return errors.New("parameters can't be zero-valued")
	}

	ctrl.SetWebServer(ws)

	ws.mux.HandleFunc("/ws/"+ctrlVer+"/"+ctrlPattern+"/", func(w http.ResponseWriter, r *http.Request) {

		ws.invokeCtrlMethod(ctrl, w, r)

	})

	return nil

}

func (ws *WebServer) invokeCtrlMethod(ctrl BaseControllerInterface, w http.ResponseWriter, r *http.Request) {

	urlElements := strings.Split(r.RequestURI, "/")

	numElements := len(urlElements)

	var args []string

	if numElements > 4 {
		args = urlElements[5:]
	} else if numElements == 4 {
		args = make([]string, 0)
	} else {
		JSONError(w, errors.New("invalid number of params"))
		return
	}

	methodName := strings.Title(urlElements[4])

	ctrl.SetRequest(r)
	ctrl.SetResponse(w)

	method := reflect.ValueOf(ctrl).MethodByName(methodName)

	if !method.IsValid() {
		JSONError(w, errors.New("method not found"))
		return
	}

	callable := method.Interface().(func([]string))

	callable(args)

}

func (ws *WebServer) GetMicroServiceRootDir() string {
	return microServiceRootDir
}
