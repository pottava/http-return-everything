package controllers

import (
	"context"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/models"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/google"
	"github.com/pottava/http-return-everything/app/lib"
	appModels "github.com/pottava/http-return-everything/app/models"
)

func getGoogleCloud(params google.GetGoogleCloudParams) middleware.Responder {
	info, found := getGoogleCloudInformation(params.HTTPRequest.Context())
	if !found {
		code := http.StatusNotFound
		return google.NewGetGoogleCloudDefault(code).WithPayload(&models.Error{
			Code:    swag.String(strconv.FormatInt(int64(code), 10)),
			Message: swag.String(http.StatusText(code)),
		})
	}
	return google.NewGetGoogleCloudOK().WithPayload(&info)
}

func getGoogleComputeEngine(params google.GetGoogleComputeEngineParams) middleware.Responder {
	meta, found := gceInstanceMetadata(params.HTTPRequest.Context())
	if !found {
		code := http.StatusNotFound
		return google.NewGetGoogleComputeEngineDefault(code).WithPayload(&models.Error{
			Code:    swag.String(strconv.FormatInt(int64(code), 10)),
			Message: swag.String(http.StatusText(code)),
		})
	}
	return google.NewGetGoogleComputeEngineOK().WithPayload(meta)
}

func getGoogleComputeEngineField(params google.GetGoogleComputeEngineFieldParams) middleware.Responder {
	code := http.StatusNotFound
	notfound := google.NewGetGoogleComputeEngineFieldDefault(code).WithPayload(&models.Error{
		Code:    swag.String(strconv.FormatInt(int64(code), 10)),
		Message: swag.String(http.StatusText(code)),
	})
	meta, found := gceInstanceMetadata(params.HTTPRequest.Context())
	if !found {
		return notfound
	}
	for _, field := range appModels.GceMetas {
		if strings.EqualFold(field.Key, params.Key) {
			return google.NewGetGoogleComputeEngineFieldOK().WithPayload(field.GCE(meta, ""))
		}
	}
	return notfound
}

func getGoogleKubernetesEngine(params google.GetGoogleKubernetesEngineParams) middleware.Responder {
	meta, found := gkeInstanceMetadata(params.HTTPRequest.Context())
	if !found {
		code := http.StatusNotFound
		return google.NewGetGoogleKubernetesEngineDefault(code).WithPayload(&models.Error{
			Code:    swag.String(strconv.FormatInt(int64(code), 10)),
			Message: swag.String(http.StatusText(code)),
		})
	}
	return google.NewGetGoogleKubernetesEngineOK().WithPayload(meta)
}

func getGoogleKubernetesEngineField(params google.GetGoogleKubernetesEngineFieldParams) middleware.Responder {
	code := http.StatusNotFound
	notfound := google.NewGetGoogleKubernetesEngineFieldDefault(code).WithPayload(&models.Error{
		Code:    swag.String(strconv.FormatInt(int64(code), 10)),
		Message: swag.String(http.StatusText(code)),
	})
	meta, found := gkeInstanceMetadata(params.HTTPRequest.Context())
	if !found {
		return notfound
	}
	for _, field := range appModels.GkeMetas {
		if strings.EqualFold(field.Key, params.Key) {
			return google.NewGetGoogleKubernetesEngineFieldOK().WithPayload(field.GKE(meta, ""))
		}
	}
	return notfound
}

func getGoogleCloudInformation(ctx context.Context) (models.GoogleCloud, bool) {
	google := models.GoogleCloud{}
	if gce, found := gceInstanceMetadata(ctx); found {
		google.Gce = gce
	}
	if gke, found := gkeInstanceMetadata(ctx); found {
		google.Gke = gke
	}
	return google, !reflect.DeepEqual(google, models.GoogleCloud{})
}

const (
	gceInstanceMetadataPrefix = "http://metadata.google.internal/computeMetadata/v1"
)

var (
	googleFlavor = &http.Header{"Metadata-Flavor": []string{"Google"}}
)

func gceInstanceMetadata(ctx context.Context) (*models.GoogleComputeEngine, bool) {
	gce := &models.GoogleComputeEngine{}
	if !lib.Config.EnabledGCP {
		return gce, false
	}
	client := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: len(appModels.GceMetas)},
		Timeout:   time.Duration(250) * time.Millisecond,
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(appModels.GceMetas))

	for _, meta := range appModels.GceMetas {
		go func(meta *appModels.GoogleCloudMetadata) {
			defer wg.Done()

			url := gceInstanceMetadataPrefix + meta.URI
			if body, err := lib.HTTPGet(ctx, client, url, googleFlavor); err == nil {
				meta.GCE(gce, string(body))
			}
		}(meta)
	}
	wg.Wait()
	return gce, !reflect.DeepEqual(*gce, models.GoogleComputeEngine{})
}

func gkeInstanceMetadata(ctx context.Context) (*models.GoogleKubernetesEngine, bool) {
	gke := &models.GoogleKubernetesEngine{}
	if !lib.Config.EnabledGCP {
		return gke, false
	}
	client := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: len(appModels.GkeMetas)},
		Timeout:   time.Duration(250) * time.Millisecond,
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(appModels.GkeMetas))

	for _, meta := range appModels.GkeMetas {
		go func(meta *appModels.GoogleCloudMetadata) {
			defer wg.Done()

			url := gceInstanceMetadataPrefix + meta.URI
			if body, err := lib.HTTPGet(ctx, client, url, googleFlavor); err == nil {
				meta.GKE(gke, string(body))
			}
		}(meta)
	}
	wg.Wait()
	return gke, !reflect.DeepEqual(*gke, models.GoogleKubernetesEngine{})
}
