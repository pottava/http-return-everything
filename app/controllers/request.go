package controllers

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/request"
)

func getRequestInfo(params request.GetRequestInfoParams) middleware.Responder {
	return request.NewGetRequestInfoOK().WithPayload(everything(params.HTTPRequest).Req)
}

func getRequestField(params request.GetRequestFieldParams) middleware.Responder {
	obj := everything(params.HTTPRequest).Req
	result := ""
	switch params.Key {
	case "protocol":
		result = swag.StringValue(obj.Protocol)
	case "method":
		result = swag.StringValue(obj.Method)
	case "host":
		result = swag.StringValue(obj.Host)
	case "remote_addr":
		result = swag.StringValue(obj.RemoteAddr)
	case "uri":
		result = swag.StringValue(obj.URI)
	}
	return request.NewGetRequestFieldOK().WithPayload(result)
}

func getRequestHeaders(params request.GetRequestHeadersParams) middleware.Responder {
	return request.NewGetRequestHeadersOK().WithPayload(everything(params.HTTPRequest).Req.Headers)
}

func getRequestHeader(params request.GetRequestHeaderParams) middleware.Responder {
	result := []string{}
	for key, value := range everything(params.HTTPRequest).Req.Headers {
		if strings.EqualFold(key, params.Header) {
			result = value
		}
	}
	return request.NewGetRequestHeaderOK().WithPayload(result)
}

func getRequestForm(params request.GetRequestFormParams) middleware.Responder {
	return request.NewGetRequestFormOK().WithPayload(everything(params.HTTPRequest).Req.Form)
}

func getRequestPostForm(params request.GetRequestPostFormParams) middleware.Responder {
	return request.NewGetRequestPostFormOK().WithPayload(everything(params.HTTPRequest).Req.PostForm)
}
