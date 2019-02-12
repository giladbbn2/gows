package gows

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	microServiceRootDir string
	templateDir         string
	serveMuxHandler     *http.ServeMux
	serveMuxHandlerOnce sync.Once
	mysqlConnections    map[string]*MysqlConnConfig
)

func init() {
	microServiceRootDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	templateDir = microServiceRootDir + string(os.PathSeparator) + "tpl"
	mysqlConnections = make(map[string]*MysqlConnConfig)
}

func Sanitize(str string) string {

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

func JSON(w http.ResponseWriter, str string) {

	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"code":200,"value":"`+str+`"}`)

}

func JSONError(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"code":500,"value":"`+err.Error()+`"}`)

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
