package models

import (
	"fmt"
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

// ToAmazonECSformat itself as AmazonECS
func (e *ECSTaskMeta) ToAmazonECS() *models.AmazonECS {
	container := ECSTaskMetaContainer{}
	if len(e.Containers) > 0 {
		container = e.Containers[0]
	}
	mappings := []*models.AmazonECSPortMappingsItems0{}
	for _, mapping := range container.Ports {
		mappings = append(mappings, &models.AmazonECSPortMappingsItems0{
			ContainerPort: swag.String(fmt.Sprintf("%d", mapping.ContainerPort)),
			HostPort:      swag.String(fmt.Sprintf("%d", mapping.HostPort)),
			Protocol:      swag.String(mapping.Protocol),
		})
	}
	networks := []*models.AmazonECSNetworksItems0{}
	for _, network := range container.Networks {
		networks = append(networks, &models.AmazonECSNetworksItems0{
			NetworkMode:   swag.String(network.NetworkMode),
			IPV4Addresses: swag.String(network.IPv4Addresses[0]),
		})
	}
	return &models.AmazonECS{
		Cluster:              e.Cluster,
		ContainerInstanceArn: nil,
		TaskArn:              swag.String(e.TaskARN),
		ContainerID:          container.ID,
		ContainerName:        swag.String(container.Name),
		DockerContainerName:  container.DockerName,
		AvailabilityZone:     e.AvailabilityZone,
		ImageID:              container.ImageID,
		ImageName:            container.Image,
		PortMappings:         mappings,
		Networks:             networks,
	}
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
	CPU    uint `json:"CPU"`
	Memory uint `json:"Memory"`
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

// ToAmazonECSformat itself as AmazonECS
func (e *ECSMetadataV1) ToAmazonECS() *models.AmazonECS {
	mappings := []*models.AmazonECSPortMappingsItems0{}
	for _, mapping := range e.PortMappings {
		mappings = append(mappings, &models.AmazonECSPortMappingsItems0{
			ContainerPort: swag.String(mapping.ContainerPort),
			HostPort:      swag.String(mapping.HostPort),
			BindIP:        swag.String(mapping.BindIP),
			Protocol:      swag.String(mapping.Protocol),
		})
	}
	networks := []*models.AmazonECSNetworksItems0{}
	for _, network := range e.Networks {
		networks = append(networks, &models.AmazonECSNetworksItems0{
			NetworkMode:   swag.String(network.NetworkMode),
			IPV4Addresses: swag.String(network.IPV4Addresses[0]),
		})
	}
	return &models.AmazonECS{
		Cluster:              e.Cluster,
		ContainerInstanceArn: swag.String(e.ContainerInstanceArn),
		TaskArn:              swag.String(e.TaskArn),
		ContainerID:          e.ContainerID,
		ContainerName:        swag.String(e.ContainerName),
		DockerContainerName:  e.DockerContainerName,
		AvailabilityZone:     "",
		ImageID:              e.ImageID,
		ImageName:            e.ImageName,
		PortMappings:         mappings,
		Networks:             networks,
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
