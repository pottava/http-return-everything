package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/models"
	"github.com/pottava/http-return-everything/app/generated/restapi/operations/aws"
)

func getAWS(params aws.GetAWSParams) middleware.Responder {
	info, found := getAWSInformation()
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
	meta, found := ec2InstanceMetadata()
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
	meta, found := ec2InstanceMetadata()
	if !found {
		return notfound
	}
	switch params.Key {
	case "instance_id":
		return aws.NewGetAmazonEC2FieldOK().WithPayload(meta.InstanceID)
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
	}
	return notfound
}

func getAmazonECS(params aws.GetAmazonECSParams) middleware.Responder {
	meta, found := ecsContainerMetadata()
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
	meta, found := ecsContainerMetadata()
	if !found {
		return notfound
	}
	switch params.Key {
	case "container_id":
		return aws.NewGetAmazonECSFieldOK().WithPayload(meta.ContainerID)
	case "container_name":
		return aws.NewGetAmazonECSFieldOK().WithPayload(swag.StringValue(meta.ContainerName))
	case "container_instance_arn":
		return aws.NewGetAmazonECSFieldOK().WithPayload(swag.StringValue(meta.ContainerInstanceArn))
	case "docker_container_name":
		return aws.NewGetAmazonECSFieldOK().WithPayload(meta.DockerContainerName)
	case "image_id":
		return aws.NewGetAmazonECSFieldOK().WithPayload(meta.ImageID)
	case "image_name":
		return aws.NewGetAmazonECSFieldOK().WithPayload(meta.ImageName)
	case "task_arn":
		return aws.NewGetAmazonECSFieldOK().WithPayload(swag.StringValue(meta.TaskArn))
	}
	return notfound
}

const (
	ec2InstanceMetadataPrefix   = "http://169.254.169.254/latest/meta-data/"
	ecsContainerMetadataFileKey = "ECS_CONTAINER_METADATA_FILE"
)

func getAWSInformation() (models.AWS, bool) {
	aws := models.AWS{}
	if ec2, found := ec2InstanceMetadata(); found {
		aws.Ec2 = &ec2
	}
	if ecs, found := ecsContainerMetadata(); found {
		aws.Ecs = ecs
	}
	return aws, !reflect.DeepEqual(aws, models.AWS{})
}

func ec2InstanceMetadata() (models.AmazonEC2, bool) {
	keys := []string{"instance-id", "placement/availability-zone", "iam/info",
		"public-hostname", "public-ipv4", "local-hostname", "local-ipv4"}
	client := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: len(keys)},
		Timeout:   time.Duration(2) * time.Second,
	}
	ec2 := models.AmazonEC2{}
	wg := &sync.WaitGroup{}
	wg.Add(len(keys))

	for _, key := range keys {
		go func(key string) {
			defer wg.Done()

			res, err := client.Get(ec2InstanceMetadataPrefix + key)
			if err != nil {
				return
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return
			}
			switch key {
			case "instance-id":
				ec2.InstanceID = string(body)
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
			}
		}(key)
	}
	wg.Wait()
	return ec2, !reflect.DeepEqual(ec2, models.AmazonEC2{})
}

type ecsPortMappings struct {
	ContainerPort string
	HostPort      string
	BindIP        string
	Protocol      string
}
type ecsNetworks struct {
	NetworkMode   string
	IPV4Addresses []string
}
type ecsMetadata struct {
	ContainerInstanceArn string `json:"ContainerInstanceARN"`
	TaskArn              string `json:"TaskARN"`
	ContainerID          string `json:"ContainerID"`
	ContainerName        string `json:"ContainerName"`
	DockerContainerName  string `json:"DockerContainerName"`
	ImageID              string `json:"ImageID"`
	ImageName            string `json:"ImageName"`
	PortMappings         []ecsPortMappings
	Networks             []ecsNetworks
	MetadataFileStatus   string `json:"MetadataFileStatus"`
}

func ecsContainerMetadata() (*models.AmazonECS, bool) {
	if value, found := os.LookupEnv(ecsContainerMetadataFileKey); found {
		file, err := ioutil.ReadFile(value)
		if err != nil {
			return nil, false
		}
		ecs := ecsMetadata{}

		type Alias ecsMetadata
		alias := &struct {
			*Alias
			PortMappings []map[string]interface{} `json:"PortMappings"`
			Networks     []map[string]interface{} `json:"Networks"`
		}{
			Alias: (*Alias)(&ecs),
		}
		json.Unmarshal(file, alias)

		var mappings []ecsPortMappings
		for _, mapping := range alias.PortMappings {
			mappings = append(mappings, ecsPortMappings{
				ContainerPort: fmt.Sprintf("%v", mapping["ContainerPort"]),
				HostPort:      fmt.Sprintf("%v", mapping["HostPort"]),
				BindIP:        fmt.Sprintf("%v", mapping["BindIp"]),
				Protocol:      fmt.Sprintf("%v", mapping["Protocol"]),
			})
		}
		ecs.PortMappings = mappings

		var networks []ecsNetworks
		for _, network := range alias.Networks {
			networks = append(networks, ecsNetworks{
				NetworkMode:   fmt.Sprintf("%v", network["NetworkMode"]),
				IPV4Addresses: interfaceToSlice(network["IPv4Addresses"]),
			})
		}
		ecs.Networks = networks

		return ecs.toAmazonECS(), true
	}
	return nil, false
}

func (e *ecsMetadata) toAmazonECS() *models.AmazonECS {
	mappings := []*models.AmazonECSPortMappingsItems0{}
	for _, mapping := range e.PortMappings {
		mappings = append(mappings, mapping.toAmazonECSPortMappingsItem())
	}
	networks := []*models.AmazonECSNetworksItems0{}
	for _, network := range e.Networks {
		networks = append(networks, network.toAmazonECSNetworksItem())
	}
	return &models.AmazonECS{
		ContainerInstanceArn: swag.String(e.ContainerInstanceArn),
		TaskArn:              swag.String(e.TaskArn),
		ContainerID:          e.ContainerID,
		ContainerName:        swag.String(e.ContainerName),
		DockerContainerName:  e.DockerContainerName,
		ImageID:              e.ImageID,
		ImageName:            e.ImageName,
		PortMappings:         mappings,
		Networks:             networks,
	}
}

func (e *ecsPortMappings) toAmazonECSPortMappingsItem() *models.AmazonECSPortMappingsItems0 {
	return &models.AmazonECSPortMappingsItems0{
		ContainerPort: swag.String(e.ContainerPort),
		HostPort:      swag.String(e.HostPort),
		BindIP:        swag.String(e.BindIP),
		Protocol:      swag.String(e.Protocol),
	}
}

func (e *ecsNetworks) toAmazonECSNetworksItem() *models.AmazonECSNetworksItems0 {
	return &models.AmazonECSNetworksItems0{
		NetworkMode:   swag.String(e.NetworkMode),
		IPV4Addresses: swag.String(e.IPV4Addresses[0]),
	}
}

func interfaceToSlice(slice interface{}) []string {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}
	ret := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = fmt.Sprintf("%v", s.Index(i).Interface())
	}
	return ret
}
