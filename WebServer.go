package gows

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
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

	fs := http.StripPrefix("/includes", http.FileServer(http.Dir("./includes")))
	mux.Handle("/includes/", fs)

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "pong")
	})

	ws := new(WebServer)
	ws.mux = mux
	ws.addrWithPort = addrWithPort

	return ws, nil

}

func (ws *WebServer) RegisterRoute(pattern string, path string) {

	ws.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		r.URL.Path = path
		ws.mux.ServeHTTP(w, r)

	})

}

func (ws *WebServer) ListenAndServe() error {

	err := http.ListenAndServe(ws.addrWithPort, ws.mux)
	if err != nil {
		return err
	}

	return nil

}

func (ws *WebServer) ListenAndServeTLS(CertPath string, PrivateKeyPath string) error {

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:         ws.addrWithPort,
		Handler:      ws.mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	err := srv.ListenAndServeTLS(CertPath, PrivateKeyPath)
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

	urlElements := strings.Split(r.URL.Path, "/")

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
