package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pottava/http-return-everything/lib"
)

var (
	version string
	date    string
)

func main() {
	http.Handle("/", lib.Wrap(index))

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if len(version) > 0 && len(date) > 0 {
			fmt.Fprintf(w, "version: %s (built at %s)", version, date)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	log.Printf("[service] listening on port %s", lib.Config.Port)
	log.Fatal(http.ListenAndServe(":"+lib.Config.Port, nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{}

	req := map[string]interface{}{}
	req["Proto"] = r.Proto
	req["Method"] = r.Method
	req["Host"] = r.Host
	req["RemoteAddr"] = r.RemoteAddr
	req["URI"] = r.RequestURI
	req["Headers"] = r.Header
	req["Form"] = r.Form
	req["PostForm"] = r.PostForm
	response["Request"] = req

	host := map[string]interface{}{}
	host["Name"], _ = os.Hostname()
	host["EnvVars"] = os.Environ()
	if data, err := ioutil.ReadFile("/etc/hosts"); err == nil {
		host["EtcHosts"] = strings.Split(string(data), "\n")
	}
	if data, err := ioutil.ReadFile("/etc/resolv.conf"); err == nil {
		host["ResolvConf"] = strings.Split(string(data), "\n")
	}
	response["Host"] = host

	app := map[string]interface{}{}
	app["Arguments"] = os.Args
	app["WorkDir"], _ = os.Getwd()
	app["Group"] = os.Getgid()
	app["User"] = os.Getuid()
	response["Application"] = app

	lib.RenderJSON(w, response, nil)
}
