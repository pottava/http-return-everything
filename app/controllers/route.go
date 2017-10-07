// Package controllers defines application's routes.
package controllers

import (
	"github.com/pottava/http-return-everything/app/generated/restapi/operations"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/application"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/host"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/request"
)

// Routes set API handlers
func Routes(api *operations.ReturnEverythingAPI) {
	api.GetEverythingHandler = operations.GetEverythingHandlerFunc(getEverything)

	api.ApplicationGetAppHandler = application.GetAppHandlerFunc(getApp)
	api.ApplicationGetAppFieldHandler = application.GetAppFieldHandlerFunc(getAppField)
	api.ApplicationGetAppEnvsHandler = application.GetAppEnvsHandlerFunc(getAppEnvs)
	api.ApplicationGetAppEnvHandler = application.GetAppEnvHandlerFunc(getAppEnv)

	api.HostGetHostHandler = host.GetHostHandlerFunc(getHost)
	api.HostGetHostFieldHandler = host.GetHostFieldHandlerFunc(getHostField)

	api.RequestGetRequestInfoHandler = request.GetRequestInfoHandlerFunc(getRequestInfo)
	api.RequestGetRequestFieldHandler = request.GetRequestFieldHandlerFunc(getRequestField)
	api.RequestGetRequestHeadersHandler = request.GetRequestHeadersHandlerFunc(getRequestHeaders)
	api.RequestGetRequestHeaderHandler = request.GetRequestHeaderHandlerFunc(getRequestHeader)
	api.RequestGetRequestFormHandler = request.GetRequestFormHandlerFunc(getRequestForm)
	api.RequestGetRequestPostFormHandler = request.GetRequestPostFormHandlerFunc(getRequestPostForm)
}
