package instances

import (
	"fmt"
	"strings"
)

type InstanceState string

const (
	InstanceStatePending      InstanceState = "pending"
	InstanceStateRunning      InstanceState = "running"
	InstanceStateShuttingDown InstanceState = "shutting-down"
	InstanceStateStopping     InstanceState = "stopping"
	InstanceStateStopped      InstanceState = "stopped"
	InstanceStateTerminated   InstanceState = "terminated"
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
