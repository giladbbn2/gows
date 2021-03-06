package gows

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync/atomic"
)

type WebServer struct {
	mux               *http.ServeMux
	addrWithPort      string
	inProgressCounter int64
}

func NewWebServer(addrWithPort string) (*WebServer, error) {

	if addrWithPort == "" {
		return nil, errors.New("parameter can't be zero-valued")
	}

	mux := http.NewServeMux()

	ws := new(WebServer)
	ws.mux = mux
	ws.addrWithPort = addrWithPort

	fs := http.StripPrefix("/includes", http.FileServer(http.Dir("./includes")))
	mux.Handle("/includes/", fs)

	mux.HandleFunc("/gows/stats", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, `{"inProgressCounter":`+strconv.FormatInt(ws.GetInProgressCounter(), 10)+`}`)
	})

	return ws, nil

}

func (ws *WebServer) RegisterRedirect(pattern string, url string, code int) {

	ws.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		http.Redirect(w, r, url, code)

	})

}

func (ws *WebServer) RegisterLocalRoute(pattern string, path string) {

	ws.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		r.URL.Path = path
		ws.mux.ServeHTTP(w, r)

	})

}

func (ws *WebServer) RegisterRemoteRoute(pattern string, remoteUrl string) {

	target, _ := url.Parse(remoteUrl) //"http://localhost:9000/"

	targetQuery := target.RawQuery

	director := func(req *http.Request) {

		req.URL.Scheme = "http"

		req.URL.Host = target.Host

		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}

		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", target.Host)

	}

	proxy := &httputil.ReverseProxy{Director: director}

	ws.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
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

		atomic.AddInt64(&ws.inProgressCounter, 1)

		invokeCtrlMethod(ctrl, w, r)

		atomic.AddInt64(&ws.inProgressCounter, -1)

	})

	return nil

}

func (ws *WebServer) GetMicroServiceRootDir() string {

	return microServiceRootDir

}

func (ws *WebServer) GetInProgressCounter() int64 {

	return atomic.LoadInt64(&ws.inProgressCounter)

}
