package gows

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

var (
	microServiceRootDir string
	serveMuxHandler     *http.ServeMux
	serveMuxHandlerOnce sync.Once
	mysqlConnections    map[string]*MysqlConnConfig
)

func init() {
	microServiceRootDir, _ = filepath.Abs(filepath.Dir(os.Args[0])) //string(os.PathSeparator)
	mysqlConnections = make(map[string]*MysqlConnConfig)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func Sanitize(str string) string {

	if str == "" {
		return ""
	}

	r := strings.NewReplacer(
		"\"", "",
		"'", "",
		"<", "",
		">", "",
		"\r", "",
		"\b", "",
		"\t", "",
		"\\", "",
		"\x00", "",
		"\n", "",
		"\x1a", "")
	return r.Replace(str)

}

func invokeCtrlMethod(ctrl BaseControllerInterface, w http.ResponseWriter, r *http.Request) {

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

func JSON(w http.ResponseWriter, str string) {

	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"code":200,"value":`+str+`}`)

}

func JSONError(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"code":500,"value":"`+err.Error()+`"}`)

}

func Gzip(b []byte) ([]byte, error) {

	if b == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write(b)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}

func Gunzip(b []byte) ([]byte, error) {

	if b == nil {
		return nil, nil
	}

	buf := bytes.NewBuffer(b)

	zr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	b2, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil, err
	}

	if err := zr.Close(); err != nil {
		return nil, err
	}

	return b2, nil

}
