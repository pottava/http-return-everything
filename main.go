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
	http.Handle("/envs/", lib.Wrap(envs))
	http.Handle("/request/", lib.Wrap(request))

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if len(version) > 0 && len(date) > 0 {
			fmt.Fprintf(w, "version: %s (built at %s)\n", version, date)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	log.Printf("[service] listening on port %s", lib.Config.Port)
	log.Fatal(http.ListenAndServe(":"+lib.Config.Port, nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{}
	response["Request"] = req(r)

	host := map[string]interface{}{}
	host["Name"], _ = os.Hostname()
	if data, err := ioutil.ReadFile("/etc/hosts"); err == nil {
		host["EtcHosts"] = strings.Split(string(data), "\n")
	}
	if data, err := ioutil.ReadFile("/etc/resolv.conf"); err == nil {
		host["ResolvConf"] = strings.Split(string(data), "\n")
	}
	response["Host"] = host

	app := map[string]interface{}{}
	app["Arguments"] = os.Args
	app["EnvVars"] = os.Environ()
	app["WorkDir"], _ = os.Getwd()
	app["Group"] = os.Getgid()
	app["User"] = os.Getuid()
	response["Application"] = app

	lib.RenderJSON(w, response, nil)
}

func req(r *http.Request) (req map[string]interface{}) {
	req = map[string]interface{}{}
	req["Proto"] = r.Proto
	req["Method"] = r.Method
	req["Host"] = r.Host
	req["RemoteAddr"] = r.RemoteAddr
	req["URI"] = r.RequestURI
	req["Headers"] = r.Header
	req["Form"] = r.Form
	req["PostForm"] = r.PostForm
	return
}

func envs(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/envs/"):]
	if len(key) == 0 {
		lib.RenderJSON(w, os.Environ(), nil)
		return
	}
	fmt.Fprintln(w, os.Getenv(key))
}

func request(w http.ResponseWriter, r *http.Request) {
	key := strings.ToLower(r.URL.Path[len("/request/"):])
	switch {
	case strings.HasPrefix(key, "proto"):
		fmt.Fprintln(w, r.Proto)
	case strings.HasPrefix(key, "method"):
		fmt.Fprintln(w, r.Method)
	case strings.HasPrefix(key, "host"):
		fmt.Fprintln(w, r.Host)
	case strings.HasPrefix(key, "remoteaddr"):
		fallthrough
	case strings.HasPrefix(key, "addr"):
		fallthrough
	case strings.HasPrefix(key, "address"):
		fmt.Fprintln(w, r.RemoteAddr)
	case strings.HasPrefix(key, "requesturi"):
		fallthrough
	case strings.HasPrefix(key, "uri"):
		fallthrough
	case strings.HasPrefix(key, "url"):
		fmt.Fprintln(w, r.RequestURI)
	case strings.HasPrefix(key, "headers"):
		if len(key) <= len("headers/") {
			lib.RenderJSON(w, r.Header, nil)
			return
		}
		if values, found := lib.Header(r, key[len("headers/"):]); found {
			lib.RenderJSON(w, values, nil)
			return
		}
		http.Error(w, "", http.StatusNotFound)
	case strings.HasPrefix(key, "form"):
		lib.RenderJSON(w, r.Form, nil)
	case strings.HasPrefix(key, "postform"):
		lib.RenderJSON(w, r.PostForm, nil)
	default:
		lib.RenderJSON(w, req(r), nil)
	}
}
