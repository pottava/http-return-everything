package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/models"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations"
)

func getEverything(params operations.GetEverythingParams) middleware.Responder {
	return operations.NewGetEverythingOK().WithPayload(everything(params.HTTPRequest))
}

func everything(r *http.Request) *models.Everything {
	envs := os.Environ()
	hosts := []string{}
	resolv := []string{}
	var aws *models.AWS
	var gcp *models.GoogleCloud

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		if data, err := ioutil.ReadFile("/etc/hosts"); err == nil {
			for _, candidate := range strings.Split(string(data), "\n") {
				if candidate != "" {
					hosts = append(hosts, candidate)
				}
			}
		}
		if data, err := ioutil.ReadFile("/etc/resolv.conf"); err == nil {
			for _, candidate := range strings.Split(string(data), "\n") {
				if candidate != "" {
					resolv = append(resolv, candidate)
				}
			}
		}
		sort.Slice(envs, func(i, j int) bool {
			return strings.ToLower(envs[i]) < strings.ToLower(envs[j])
		})
	}()
	go func() {
		defer wg.Done()
		if candidate, found := getAWSInformation(r.Context()); found {
			aws = &candidate
		}
	}()
	go func() {
		defer wg.Done()
		if candidate, found := getGoogleCloudInformation(r.Context()); found {
			gcp = &candidate
		}
	}()
	wg.Wait()

	host, _ := os.Hostname()
	wd, _ := os.Getwd()
	r.ParseForm()

	return &models.Everything{
		App: &models.Application{
			Args:    os.Args,
			Envs:    envs,
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
		Aws:         aws,
		Googlecloud: gcp,
	}
}
