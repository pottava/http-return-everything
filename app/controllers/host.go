package controllers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/host"
)

func getHost(params host.GetHostParams) middleware.Responder {
	return host.NewGetHostOK().WithPayload(everything(params.HTTPRequest).Host)
}

func getHostField(params host.GetHostFieldParams) middleware.Responder {
	obj := everything(params.HTTPRequest).Host
	var result interface{}
	switch params.Key {
	case "name":
		result = swag.StringValue(obj.Name)
	case "hosts":
		result = obj.Hosts
	case "resolv_conf":
		result = obj.ResolvConf
	}
	return host.NewGetHostFieldOK().WithPayload(result)
}
