// Package controllers defines application's routes.
package controllers

import (
	"github.com/pottava/http-return-everything/app/generated/restapi/operations"
)

// Routes set API handlers
func Routes(api *operations.ReturnEverythingAPI) {
	api.GetEverythingHandler = operations.GetEverythingHandlerFunc(getEverything)
	api.GetAppHandler = operations.GetAppHandlerFunc(getApp)
	api.GetAppEnvsHandler = operations.GetAppEnvsHandlerFunc(getAppEnvs)
	api.GetHostHandler = operations.GetHostHandlerFunc(getHost)
	api.GetRequestInfoHandler = operations.GetRequestInfoHandlerFunc(getRequestInfo)
}
