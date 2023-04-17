package instances

import "fmt"

type CloudProvider interface {
	StartInstance(id string) error
	StopInstance(id string) error
	GetInstanceStatus(id string) (InstanceState, error)
	GetName() string
}

type MockAWSCloud struct {
}

func (m MockAWSCloud) StartInstance(id string) error {
	fmt.Println("starting AWS instance ", id)
	return nil
}

func (m MockAWSCloud) StopInstance(id string) error {
	fmt.Println("stopping AWS instance ", id)
	return nil
}

func (m MockAWSCloud) GetInstanceStatus(id string) (InstanceState, error) {
	switch id {
	case "NotExist":
		return "", fmt.Errorf("instance %q not found in the cloud provider", id)
	case "Running":
		return InstanceStateRunning, nil
	case "Stopped":
		return InstanceStateStopped, nil

	default:
		return InstanceStateRunning, nil

	}
}

func (m MockAWSCloud) GetName() string {
	return "mock-aws"
}
