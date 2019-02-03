package gows

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type BaseController struct {
	MicroServiceRootDir string
	TemplateDir         string
	Request             *http.Request
	Response            http.ResponseWriter
}

func (ctrl *BaseController) SetRequest(r *http.Request) {
	ctrl.Request = r
}

func (ctrl *BaseController) SetResponse(w http.ResponseWriter) {
	ctrl.Response = w
}

func (ctrl *BaseController) SetDirs() {
	ctrl.MicroServiceRootDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	ctrl.TemplateDir = ctrl.MicroServiceRootDir + string(os.PathSeparator) + "tpl"
}

func (ctrl *BaseController) GetTemplate(tpl string) ([]byte, error) {

	contents, err := ioutil.ReadFile(ctrl.TemplateDir + string(os.PathSeparator) + tpl)

	if err != nil {
		return nil, err
	}

	return contents, nil

}
