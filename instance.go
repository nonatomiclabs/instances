package instances

import (
	"fmt"
	"strings"
)

type InstanceState string

const (
	InstanceStateRunning InstanceState = "running"
	InstanceStateStopped InstanceState = "stopped"
)

type Instance struct {
	Id                string `json:"id"`
	CloudProviderName string `json:"cloud-provider"`
}

func (i Instance) GetCloudProvider(cloudProviders map[string]CloudProvider) (CloudProvider, error) {
	cloudProvider, exists := cloudProviders[strings.ToLower(i.CloudProviderName)]
	if !exists {
		return nil, fmt.Errorf("unsupported cloud provider %q", i.CloudProviderName)
	}
	return cloudProvider, nil
}