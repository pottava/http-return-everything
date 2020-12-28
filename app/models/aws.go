package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/models"
)

// ECSTaskMeta is a model for ECS task metadata
type ECSTaskMeta struct {
	Cluster          string                 `json:"Cluster"`
	TaskARN          string                 `json:"TaskARN"`
	Family           string                 `json:"Family"`
	Revision         string                 `json:"Revision"`
	DesiredStatus    string                 `json:"DesiredStatus,omitempty"`
	KnownStatus      string                 `json:"KnownStatus"`
	Containers       []ECSTaskMetaContainer `json:"Containers,omitempty"`
	Limits           ECSTaskMetaLimits      `json:"Limits,omitempty"`
	PullStartedAt    *time.Time             `json:"PullStartedAt"`
	PullStoppedAt    *time.Time             `json:"PullStoppedAt"`
	AvailabilityZone string                 `json:"AvailabilityZone"`
}

// ToAmazonECS itself as AmazonECS
func (e *ECSTaskMeta) ToAmazonECS() *models.AmazonECS {
	containers := []*models.AmazonECSContainer{}
	for _, container := range e.Containers {
		mappings := []*models.AmazonECSContainerPortMappingsItems0{}
		for _, mapping := range container.Ports {
			mappings = append(mappings, &models.AmazonECSContainerPortMappingsItems0{
				ContainerPort: swag.String(fmt.Sprintf("%d", mapping.ContainerPort)),
				HostPort:      swag.String(fmt.Sprintf("%d", mapping.HostPort)),
				Protocol:      swag.String(mapping.Protocol),
			})
		}
		networks := []*models.AmazonECSContainerNetworksItems0{}
		for _, network := range container.Networks {
			networks = append(networks, &models.AmazonECSContainerNetworksItems0{
				NetworkMode:   swag.String(network.NetworkMode),
				IPV4Addresses: swag.String(strings.Join(network.IPv4Addresses, ",")),
				IPV6Addresses: strings.Join(network.IPv6Addresses, ","),
			})
		}
		containers = append(containers, &models.AmazonECSContainer{
			ID:           container.ID,
			Name:         swag.String(container.Name),
			Type:         container.Type,
			DockerName:   container.DockerName,
			ImageID:      container.ImageID,
			ImageName:    swag.String(container.Image),
			PortMappings: mappings,
			Networks:     networks,
			Desired:      container.DesiredStatus,
			Known:        container.KnownStatus,
			CPU:          fmt.Sprintf("%v", container.Limits.CPU),
			Memory:       fmt.Sprintf("%v", container.Limits.Memory),
			CreatedAt:    timeToStr(container.CreatedAt),
			StartedAt:    timeToStr(container.StartedAt),
		})
	}
	return &models.AmazonECS{
		Cluster:          swag.String(e.Cluster),
		TaskArn:          swag.String(e.TaskARN),
		Family:           e.Family,
		Revision:         e.Revision,
		Containers:       containers,
		Desired:          e.DesiredStatus,
		Known:            e.KnownStatus,
		AvailabilityZone: e.AvailabilityZone,
		CPU:              fmt.Sprintf("%v", e.Limits.CPU),
		Memory:           fmt.Sprintf("%v", e.Limits.Memory),
		PullStartedAt:    timeToStr(e.PullStartedAt),
	}
}

func timeToStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// ECSTaskMetaContainer is a model for ECS task container
type ECSTaskMetaContainer struct {
	ID            string               `json:"DockerId"`
	Name          string               `json:"Name"`
	DockerName    string               `json:"DockerName"`
	Image         string               `json:"Image"`
	ImageID       string               `json:"ImageID"`
	Ports         []ECSTaskMetaPort    `json:"Ports,omitempty"`
	Labels        map[string]string    `json:"Labels,omitempty"`
	DesiredStatus string               `json:"DesiredStatus"`
	KnownStatus   string               `json:"KnownStatus"`
	ExitCode      *int                 `json:"ExitCode,omitempty"`
	Limits        ECSTaskMetaLimits    `json:"Limits,omitempty"`
	CreatedAt     *time.Time           `json:"CreatedAt,omitempty"`
	StartedAt     *time.Time           `json:"StartedAt,omitempty"`
	FinishedAt    *time.Time           `json:"FinishedAt,omitempty"`
	Type          string               `json:"Type,omitempty"`
	Networks      []ECSTaskMetaNetwork `json:"Networks,omitempty"`
}

// ECSTaskMetaLimits is a model for ECS task limits
type ECSTaskMetaLimits struct {
	CPU    float64 `json:"CPU"`
	Memory float64 `json:"Memory"`
}

// ECSTaskMetaPort is a model for ECS task port
type ECSTaskMetaPort struct {
	Protocol      string `json:"Protocol"`
	ContainerPort uint16 `json:"ContainerPort"`
	HostPort      uint16 `json:"HostPort,omitempty"`
}

// ECSTaskMetaNetwork is a model for ECS task network
type ECSTaskMetaNetwork struct {
	NetworkMode   string   `json:"NetworkMode,omitempty"`
	IPv4Addresses []string `json:"IPv4Addresses,omitempty"`
	IPv6Addresses []string `json:"IPv6Addresses,omitempty"`
}

// ECSMetadataV1 is a model for meta API
type ECSMetadataV1 struct {
	Cluster              string `json:"Cluster"`
	ContainerInstanceArn string `json:"ContainerInstanceARN"`
	TaskArn              string `json:"TaskARN"`
	ContainerID          string `json:"ContainerID"`
	ContainerName        string `json:"ContainerName"`
	DockerContainerName  string `json:"DockerContainerName"`
	ImageID              string `json:"ImageID"`
	ImageName            string `json:"ImageName"`
	PortMappings         []ECSPortMappingsV1
	Networks             []ECSNetworksV1
	MetadataFileStatus   string `json:"MetadataFileStatus"`
}

// ToAmazonECS itself as AmazonECS
func (e *ECSMetadataV1) ToAmazonECS() *models.AmazonECS {
	containers := []*models.AmazonECSContainer{}
	mappings := []*models.AmazonECSContainerPortMappingsItems0{}
	for _, mapping := range e.PortMappings {
		mappings = append(mappings, &models.AmazonECSContainerPortMappingsItems0{
			ContainerPort: swag.String(mapping.ContainerPort),
			HostPort:      swag.String(mapping.HostPort),
			BindIP:        mapping.BindIP,
			Protocol:      swag.String(mapping.Protocol),
		})
	}
	networks := []*models.AmazonECSContainerNetworksItems0{}
	for _, network := range e.Networks {
		networks = append(networks, &models.AmazonECSContainerNetworksItems0{
			NetworkMode:   swag.String(network.NetworkMode),
			IPV4Addresses: swag.String(strings.Join(network.IPV4Addresses, ",")),
		})
	}
	containers = append(containers, &models.AmazonECSContainer{
		ID:           e.ContainerID,
		Name:         swag.String(e.ContainerName),
		DockerName:   e.DockerContainerName,
		ImageID:      e.ImageID,
		ImageName:    swag.String(e.ImageName),
		PortMappings: mappings,
		Networks:     networks,
	})
	return &models.AmazonECS{
		Cluster:              swag.String(e.Cluster),
		ContainerInstanceArn: e.ContainerInstanceArn,
		TaskArn:              swag.String(e.TaskArn),
		Containers:           containers,
	}
}

// ECSPortMappingsV1 is a model for PortMappings
type ECSPortMappingsV1 struct {
	ContainerPort string
	HostPort      string
	BindIP        string
	Protocol      string
}

// ECSNetworksV1 is a model for Networks
type ECSNetworksV1 struct {
	NetworkMode   string
	IPV4Addresses []string
}
