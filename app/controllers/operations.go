package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/models"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations"
)

func getEverything(params operations.GetEverythingParams) middleware.Responder {
	return operations.NewGetEverythingOK().WithPayload(everything(params.HTTPRequest))
}

func everything(r *http.Request) *models.Everything {
	host, _ := os.Hostname()
	wd, _ := os.Getwd()

	hosts := []string{}
	if data, err := ioutil.ReadFile("/etc/hosts"); err == nil {
		for _, candidate := range strings.Split(string(data), "\n") {
			if candidate != "" {
				hosts = append(hosts, candidate)
			}
		}
	}
	resolv := []string{}
	if data, err := ioutil.ReadFile("/etc/resolv.conf"); err == nil {
		for _, candidate := range strings.Split(string(data), "\n") {
			if candidate != "" {
				resolv = append(resolv, candidate)
			}
		}
	}
	r.ParseForm()

	return &models.Everything{
		App: &models.Application{
			Args:    os.Args,
			Envs:    os.Environ(),
			Grp:     swag.Int64(int64(os.Getgid())),
			User:    swag.Int64(int64(os.Getuid())),
			Workdir: swag.String(wd),
		},
		Host: &models.Host{
			Name:       swag.String(host),
			Hosts:      hosts,
			ResolvConf: resolv,
		},
		Req: &models.HTTPRequest{
			Protocol:   swag.String(r.Proto),
			Method:     swag.String(r.Method),
			Host:       swag.String(r.Host),
			RemoteAddr: swag.String(r.RemoteAddr),
			URI:        swag.String(r.RequestURI),
			Headers:    r.Header,
			Form:       r.Form,
			PostForm:   r.PostForm,
		},
	}
}
