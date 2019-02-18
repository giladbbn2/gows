package gows

import (
	"io/ioutil"
	"net/http"
	"os"
)

type BaseController struct {
	ws       *WebServer
	Request  *http.Request
	Response http.ResponseWriter
}

func (ctrl *BaseController) SetWebServer(ws *WebServer) {
	ctrl.ws = ws
}

func (ctrl *BaseController) GetWebServer() *WebServer {
	return ctrl.ws
}

func (ctrl *BaseController) SetRequest(r *http.Request) {
	ctrl.Request = r
}

func (ctrl *BaseController) SetResponse(w http.ResponseWriter) {
	ctrl.Response = w
}

func (ctrl *BaseController) GetResponse() *http.ResponseWriter {
	return &ctrl.Response
}

func (ctrl *BaseController) GetTemplate(tplFilename string) ([]byte, error) {

	contents, err := ioutil.ReadFile(templateDir + string(os.PathSeparator) + tplFilename)

	if err != nil {
		return nil, err
	}

	return contents, nil

}

func (ctrl *BaseController) GetMicroServiceRootDir() string {
	return microServiceRootDir
}

func (ctrl *BaseController) GetTemplateDir() string {
	return templateDir
}
