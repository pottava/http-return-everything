package controllers

import (
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/application"
)

func getApp(params application.GetAppParams) middleware.Responder {
	return application.NewGetAppOK().WithPayload(everything(params.HTTPRequest).App)
}

func getAppField(params application.GetAppFieldParams) middleware.Responder {
	app := everything(params.HTTPRequest).App
	result := ""
	switch params.Key {
	case "args":
		result = strings.Join(app.Args, " ")
	case "grp":
		result = strconv.FormatInt(swag.Int64Value(app.Grp), 10)
	case "user":
		result = strconv.FormatInt(swag.Int64Value(app.User), 10)
	case "workdir":
		result = swag.StringValue(app.Workdir)
	}
	return application.NewGetAppFieldOK().WithPayload(result)
}

func getAppEnvs(params application.GetAppEnvsParams) middleware.Responder {
	return application.NewGetAppEnvsOK().WithPayload(everything(params.HTTPRequest).App.Envs)
}

func getAppEnv(params application.GetAppEnvParams) middleware.Responder {
	result := ""
	keys := []string{
		params.Env + "=",
		strings.ToUpper(params.Env) + "=",
		strings.ToLower(params.Env) + "=",
	}
	for _, env := range everything(params.HTTPRequest).App.Envs {
		for _, key := range keys {
			if strings.HasPrefix(env, key) {
				result = strings.TrimPrefix(env, key)
			}
		}
	}
	return application.NewGetAppEnvOK().WithPayload(result)
}
