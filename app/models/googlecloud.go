package models

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/pottava/http-return-everything/app/generated/models"
)

// GoogleCloudMetadata represents metadata of Google Cloud GCE/GKE
type GoogleCloudMetadata struct {
	Key string
	URI string
}

// Variables
var (
	GceMetas  = []*GoogleCloudMetadata{}
	GkeMetas  = []*GoogleCloudMetadata{}
	tokenizer = regexp.MustCompile("[\n ]+")
)

func init() {
	// GCE
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "project_id",
		URI: "/project/project-id",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "numeric_project_id",
		URI: "/project/numeric-project-id",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "enable_os_login",
		URI: "/instance/enable-oslogin",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "instance_hostname",
		URI: "/instance/hostname",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "instance_id",
		URI: "/instance/id",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "instance_name",
		URI: "/instance/name",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "machine_type",
		URI: "/instance/machine-type",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "cpu_platform",
		URI: "/instance/cpu-platform",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "tags",
		URI: "/instance/tags",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "zone",
		URI: "/instance/zone",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "service_account_email",
		URI: "/instance/service-accounts/default/email",
	})
	GceMetas = append(GceMetas, &GoogleCloudMetadata{
		Key: "service_account_scopes",
		URI: "/instance/service-accounts/default/scopes",
	})

	// GKE
	GkeMetas = append(GkeMetas, &GoogleCloudMetadata{
		Key: "cluster_uid",
		URI: "/instance/attributes/cluster-uid",
	})
	GkeMetas = append(GkeMetas, &GoogleCloudMetadata{
		Key: "cluster_name",
		URI: "/instance/attributes/cluster-name",
	})
	GkeMetas = append(GkeMetas, &GoogleCloudMetadata{
		Key: "cluster_location",
		URI: "/instance/attributes/cluster-location",
	})
}

// GCE set value to GCE object / returns from GCE object
func (m *GoogleCloudMetadata) GCE(gce *models.GoogleComputeEngine, value string) string {
	switch m.Key {
	case "project_id":
		if len(value) == 0 {
			return swag.StringValue(gce.ProjectID)
		}
		gce.ProjectID = swag.String(value)
	case "numeric_project_id":
		if len(value) == 0 {
			return gce.NumericProjectID
		}
		gce.NumericProjectID = value
	case "enable_os_login":
		if len(value) == 0 {
			return gce.EnableOsLogin
		}
		gce.EnableOsLogin = value
	case "instance_hostname":
		if len(value) == 0 {
			return gce.InstanceHostname
		}
		gce.InstanceHostname = value
	case "instance_id":
		if len(value) == 0 {
			return swag.StringValue(gce.InstanceID)
		}
		gce.InstanceID = swag.String(value)
	case "instance_name":
		if len(value) == 0 {
			return gce.InstanceName
		}
		gce.InstanceName = value
	case "machine_type":
		if len(value) == 0 {
			return gce.MachineType
		}
		gce.MachineType = value
	case "cpu_platform":
		if len(value) == 0 {
			return gce.CPUPlatform
		}
		gce.CPUPlatform = value
	case "service_account_email":
		if len(value) == 0 {
			if len(gce.ServiceAccounts) > 0 {
				return gce.ServiceAccounts[0].Email
			}
			return ""
		}
		if len(gce.ServiceAccounts) == 0 {
			gce.ServiceAccounts = append(gce.ServiceAccounts, &models.GoogleServiceAccount{})
		}
		gce.ServiceAccounts[0].Email = value
	case "service_account_scopes":
		if len(value) == 0 {
			if len(gce.ServiceAccounts) > 0 {
				return strings.Join(gce.ServiceAccounts[0].Scopes, ",")
			}
			return ""
		}
		if len(gce.ServiceAccounts) == 0 {
			gce.ServiceAccounts = append(gce.ServiceAccounts, &models.GoogleServiceAccount{})
		}
		gce.ServiceAccounts[0].Scopes = deleteEmpty(tokenizer.Split(value, -1))
	case "tags":
		if len(value) == 0 {
			return strings.Join(gce.Tags, ",")
		}
		json.Unmarshal([]byte(value), &gce.Tags)
	case "zone":
		if len(value) == 0 {
			return gce.Zone
		}
		gce.Zone = value
	}
	return ""
}

// GKE set value to GKE object / returns from GKE object
func (m *GoogleCloudMetadata) GKE(gke *models.GoogleKubernetesEngine, value string) string {
	switch m.Key {
	case "cluster_uid":
		if len(value) == 0 {
			return gke.ClusterUID
		}
		gke.ClusterUID = value
	case "cluster_name":
		if len(value) == 0 {
			return gke.ClusterName
		}
		gke.ClusterName = value
	case "cluster_location":
		if len(value) == 0 {
			return gke.ClusterLocation
		}
		gke.ClusterLocation = value
	}
	return ""
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
