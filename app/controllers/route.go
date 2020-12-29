// Package controllers defines application's routes.
package controllers

import (
	"github.com/pottava/http-return-everything/app/generated/restapi/operations"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/application"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/aws"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/google"
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

	api.AwsGetAWSHandler = aws.GetAWSHandlerFunc(getAWS)
	api.AwsGetAmazonEC2Handler = aws.GetAmazonEC2HandlerFunc(getAmazonEC2)
	api.AwsGetAmazonEC2FieldHandler = aws.GetAmazonEC2FieldHandlerFunc(getAmazonEC2Field)
	api.AwsGetAmazonECSHandler = aws.GetAmazonECSHandlerFunc(getAmazonECS)
	api.AwsGetAmazonECSFieldHandler = aws.GetAmazonECSFieldHandlerFunc(getAmazonECSField)

	api.GoogleGetGoogleCloudHandler = google.GetGoogleCloudHandlerFunc(getGoogleCloud)
	api.GoogleGetGoogleComputeEngineHandler = google.GetGoogleComputeEngineHandlerFunc(getGoogleComputeEngine)
	api.GoogleGetGoogleComputeEngineFieldHandler = google.GetGoogleComputeEngineFieldHandlerFunc(getGoogleComputeEngineField)
	api.GoogleGetGoogleKubernetesEngineHandler = google.GetGoogleKubernetesEngineHandlerFunc(getGoogleKubernetesEngine)
	api.GoogleGetGoogleKubernetesEngineFieldHandler = google.GetGoogleKubernetesEngineFieldHandlerFunc(getGoogleKubernetesEngineField)
}
