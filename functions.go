package gows

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
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

func IP2Long(ip string) uint32 {

	return binary.BigEndian.Uint32(net.ParseIP(ip)[12:16])

}

func Long2IP(ip int64) string {

	b0 := strconv.FormatInt((ip>>24)&0xff, 10)
	b1 := strconv.FormatInt((ip>>16)&0xff, 10)
	b2 := strconv.FormatInt((ip>>8)&0xff, 10)
	b3 := strconv.FormatInt((ip & 0xff), 10)

	return b0 + "." + b1 + "." + b2 + "." + b3

}

func GetRemoteAddress(r *http.Request) (string, error) {

	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		addresses := strings.Split(ips, ",")
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || IsIPInPrivateSubnet(ip) {
				continue
			}
			return ip, nil
		}
	}

	ips = r.Header.Get("X-Real-Ip")
	if ips != "" {
		addresses := strings.Split(ips, ",")
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || IsIPInPrivateSubnet(ip) {
				continue
			}
			return ip, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	return ip, err

}

func IsIPInPrivateSubnet(ip string) bool {

	ipLong := IP2Long(ip)

	if (ipLong >= 167772160 && ipLong <= 184549375) || (ipLong >= 1681915904 && ipLong <= 1686110207) || (ipLong >= 2886729728 && ipLong <= 2887778303) || (ipLong >= 3221225472 && ipLong <= 3221225727) || (ipLong >= 3232235520 && ipLong <= 3232301055) || (ipLong >= 3323068416 && ipLong <= 3323199487) {
		return true
	}

	return false

}
