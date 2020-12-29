package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/models"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/aws"
	"github.com/pottava/http-return-everything/app/lib"
	appModels "github.com/pottava/http-return-everything/app/models"
)

func getAWS(params aws.GetAWSParams) middleware.Responder {
	info, found := getAWSInformation(params.HTTPRequest.Context())
	if !found {
		code := http.StatusNotFound
		return aws.NewGetAWSDefault(code).WithPayload(&models.Error{
			Code:    swag.String(strconv.FormatInt(int64(code), 10)),
			Message: swag.String(http.StatusText(code)),
		})
	}
	return aws.NewGetAWSOK().WithPayload(&info)
}

func getAmazonEC2(params aws.GetAmazonEC2Params) middleware.Responder {
	meta, found := ec2InstanceMetadata(params.HTTPRequest.Context())
	if !found {
		code := http.StatusNotFound
		return aws.NewGetAmazonEC2Default(code).WithPayload(&models.Error{
			Code:    swag.String(strconv.FormatInt(int64(code), 10)),
			Message: swag.String(http.StatusText(code)),
		})
	}
	return aws.NewGetAmazonEC2OK().WithPayload(&meta)
}

func getAmazonEC2Field(params aws.GetAmazonEC2FieldParams) middleware.Responder {
	code := http.StatusNotFound
	notfound := aws.NewGetAmazonEC2FieldDefault(code).WithPayload(&models.Error{
		Code:    swag.String(strconv.FormatInt(int64(code), 10)),
		Message: swag.String(http.StatusText(code)),
	})
	meta, found := ec2InstanceMetadata(params.HTTPRequest.Context())
	if !found {
		return notfound
	}
	switch params.Key {
	case "instance_id":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.InstanceID)
	case "instance_type":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.InstanceType)
	case "ami_id":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.AmiID)
	case "instance_profile":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.InstanceProfile)
	case "availability_zone":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.AvailabilityZone)
	case "public_hostname":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.PublicHostname)
	case "public_ipv4":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.PublicIPV4)
	case "local_hostname":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.LocalHostname)
	case "local_ipv4":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.LocalIPV4)
	case "security_groups":
		result := strings.Join(meta.SecurityGroups, ",")
		return aws.NewGetAmazonEC2FieldOK().WithPayload(result)
	}
	return notfound
}

func getAmazonECS(params aws.GetAmazonECSParams) middleware.Responder {
	meta, found := ecsContainerMetadata(params.HTTPRequest.Context())
	if !found {
		code := http.StatusNotFound
		return aws.NewGetAmazonECSDefault(code).WithPayload(&models.Error{
			Code:    swag.String(strconv.FormatInt(int64(code), 10)),
			Message: swag.String(http.StatusText(code)),
		})
	}
	return aws.NewGetAmazonECSOK().WithPayload(meta)
}

func getAmazonECSField(params aws.GetAmazonECSFieldParams) middleware.Responder {
	code := http.StatusNotFound
	notfound := aws.NewGetAmazonECSFieldDefault(code).WithPayload(&models.Error{
		Code:    swag.String(strconv.FormatInt(int64(code), 10)),
		Message: swag.String(http.StatusText(code)),
	})
	meta, found := ecsContainerMetadata(params.HTTPRequest.Context())
	if !found {
		return notfound
	}
	switch params.Key {
	case "cluster":
		return aws.NewGetAmazonECSFieldOK().WithPayload(swag.StringValue(meta.Cluster))
	case "container_id":
		return aws.NewGetAmazonECSFieldOK().WithPayload(meta.Containers[0].ID)
	case "container_name":
		return aws.NewGetAmazonECSFieldOK().WithPayload(swag.StringValue(meta.Containers[0].Name))
	case "container_instance_arn":
		return aws.NewGetAmazonECSFieldOK().WithPayload(meta.ContainerInstanceArn)
	case "docker_container_name":
		return aws.NewGetAmazonECSFieldOK().WithPayload(meta.Containers[0].DockerName)
	case "availability_zone":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.AvailabilityZone)
	case "image_id":
		return aws.NewGetAmazonECSFieldOK().WithPayload(meta.Containers[0].ImageID)
	case "image_name":
		return aws.NewGetAmazonECSFieldOK().WithPayload(swag.StringValue(meta.Containers[0].ImageName))
	case "task_arn":
		return aws.NewGetAmazonECSFieldOK().WithPayload(swag.StringValue(meta.TaskArn))
	}
	return notfound
}

const (
	ec2InstanceMetadataPrefix  = "http://169.254.169.254/latest/meta-data/"
	ec2InstanceAPIToken        = "http://169.254.169.254/latest/api/token" // nolint:gosec
	v4ecsTaskMetadataEndpoint  = "ECS_CONTAINER_METADATA_URI_V4"
	v3ecsTaskMetadataEndpoint  = "ECS_CONTAINER_METADATA_URI"
	v2ecsTaskMetadataEndpoint  = "http://169.254.170.2/v2/metadata"
	v1ecsContainerMetadataFile = "ECS_CONTAINER_METADATA_FILE"
)

func getAWSInformation(ctx context.Context) (models.AWS, bool) {
	aws := models.AWS{}
	if ec2, found := ec2InstanceMetadata(ctx); found {
		aws.Ec2 = &ec2
	}
	if ecs, found := ecsContainerMetadata(ctx); found {
		aws.Ecs = ecs
	}
	return aws, !reflect.DeepEqual(aws, models.AWS{})
}

func ec2InstanceMetadata(ctx context.Context) (models.AmazonEC2, bool) {
	ec2 := models.AmazonEC2{}
	if !lib.Config.EnabledAWS {
		return ec2, false
	}
	keys := []string{"instance-id", "placement/availability-zone", "iam/info",
		"public-hostname", "public-ipv4", "local-hostname", "local-ipv4"}
	client := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: len(keys)},
		Timeout:   time.Duration(250) * time.Millisecond,
	}
	token := &http.Header{"X-aws-ec2-metadata-token-ttl-seconds": []string{"3600"}}
	candidate, err := lib.HTTPPut(ctx, client, ec2InstanceAPIToken, token)
	if err != nil {
		return ec2, false
	}
	token = &http.Header{"X-aws-ec2-metadata-token": []string{string(candidate)}}

	wg := &sync.WaitGroup{}
	wg.Add(len(keys))

	for _, key := range keys {
		go func(key string) {
			defer wg.Done()

			body, err := lib.HTTPGet(ctx, client, ec2InstanceMetadataPrefix+key, token)
			if err != nil {
				return
			}
			switch key {
			case "instance-id":
				ec2.InstanceID = string(body)
			case "instance-type":
				ec2.InstanceType = string(body)
			case "ami-id":
				ec2.AmiID = string(body)
			case "placement/availability-zone":
				ec2.AvailabilityZone = string(body)
			case "iam/info":
				parse := &struct {
					InstanceProfileArn string
				}{}
				json.Unmarshal(body, parse)
				ec2.InstanceProfile = parse.InstanceProfileArn
			case "public-hostname":
				ec2.PublicHostname = string(body)
			case "public-ipv4":
				ec2.PublicIPV4 = string(body)
			case "local-hostname":
				ec2.LocalHostname = string(body)
			case "local-ipv4":
				ec2.LocalIPV4 = string(body)
			case "security-groups":
				ec2.SecurityGroups = strings.Split(string(body), "\n")
			}
		}(key)
	}
	wg.Wait()
	return ec2, !reflect.DeepEqual(ec2, models.AmazonEC2{})
}

func ecsContainerMetadata(ctx context.Context) (*models.AmazonECS, bool) {
	if !lib.Config.EnabledAWS {
		return nil, false
	}
	if meta, found := ecsContainerMetadataVx(ctx, v4ecsTaskMetadataEndpoint); found {
		return meta, found
	}
	if meta, found := ecsContainerMetadataVx(ctx, v3ecsTaskMetadataEndpoint); found {
		return meta, found
	}
	if meta, found := ecsContainerMetadataV2(ctx); found {
		return meta, found
	}
	return ecsContainerMetadataV1(ctx)
}

func ecsContainerMetadataVx(ctx context.Context, endpoint string) (*models.AmazonECS, bool) {
	if value, found := os.LookupEnv(endpoint); found {
		client := &http.Client{
			Timeout: time.Duration(250) * time.Millisecond,
		}
		body, err := lib.HTTPGet(ctx, client, value+"/task", nil)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return nil, false
		}
		meta := appModels.ECSTaskMeta{}
		if err = json.Unmarshal(body, &meta); err != nil {
			log.Printf("Error: %s", err.Error())
			return nil, false
		}
		return meta.ToAmazonECS(), true
	}
	return nil, false
}

func ecsContainerMetadataV2(ctx context.Context) (*models.AmazonECS, bool) {
	client := &http.Client{
		Timeout: time.Duration(250) * time.Millisecond,
	}
	body, err := lib.HTTPGet(ctx, client, v2ecsTaskMetadataEndpoint, nil)
	if err != nil {
		return nil, false
	}
	meta := appModels.ECSTaskMeta{}
	if err = json.Unmarshal(body, &meta); err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, false
	}
	return meta.ToAmazonECS(), true
}

func ecsContainerMetadataV1(ctx context.Context) (*models.AmazonECS, bool) {
	if value, found := os.LookupEnv(v1ecsContainerMetadataFile); found {
		file, err := ioutil.ReadFile(value)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return nil, false
		}
		type Alias appModels.ECSMetadataV1
		ecs := appModels.ECSMetadataV1{}
		alias := &struct {
			*Alias
			PortMappings []map[string]interface{} `json:"PortMappings"`
			Networks     []map[string]interface{} `json:"Networks"`
		}{
			Alias: (*Alias)(&ecs),
		}
		json.Unmarshal(file, alias)

		mappings := []appModels.ECSPortMappingsV1{}
		for _, mapping := range alias.PortMappings {
			mappings = append(mappings, appModels.ECSPortMappingsV1{
				ContainerPort: fmt.Sprintf("%v", mapping["ContainerPort"]),
				HostPort:      fmt.Sprintf("%v", mapping["HostPort"]),
				BindIP:        fmt.Sprintf("%v", mapping["BindIp"]),
				Protocol:      fmt.Sprintf("%v", mapping["Protocol"]),
			})
		}
		ecs.PortMappings = mappings

		networks := []appModels.ECSNetworksV1{}
		for _, network := range alias.Networks {
			networks = append(networks, appModels.ECSNetworksV1{
				NetworkMode:   fmt.Sprintf("%v", network["NetworkMode"]),
				IPV4Addresses: lib.InterfaceToSlice(network["IPv4Addresses"]),
			})
		}
		ecs.Networks = networks
		return ecs.ToAmazonECS(), true
	}
	return nil, false
}
