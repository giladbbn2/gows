package gows

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
)

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

func HTTPCallCtrlMethod(ctrl BaseControllerInterface, w http.ResponseWriter, r *http.Request) error {

	urlElements := strings.Split(r.RequestURI, "/")

	numElements := len(urlElements)

	var args []string

	if numElements > 4 {
		args = urlElements[5:]
	} else if numElements == 4 {
		args = make([]string, 0)
	} else {
		return errors.New("invalid number of params")
	}

	methodName := strings.Title(urlElements[4])

	ctrl.SetDirs()
	ctrl.SetRequest(r)
	ctrl.SetResponse(w)

	method := reflect.ValueOf(ctrl).MethodByName(methodName)

	if !method.IsValid() {
		return errors.New("method not found")
	}

	callable := method.Interface().(func([]string))

	callable(args)

	return nil

}
