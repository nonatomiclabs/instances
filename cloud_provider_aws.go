package instances

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2InstanceManager interface {
	DescribeInstanceStatus(ctx context.Context, params *ec2.DescribeInstanceStatusInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error)
	StartInstances(ctx context.Context, params *ec2.StartInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
	StopInstances(ctx context.Context, params *ec2.StopInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
}

type AWSCloud struct {
	Ec2Client EC2InstanceManager
}

func (a AWSCloud) StartInstance(id string) error {
	ctx := context.TODO()
	state, err := a.GetInstanceStatus(id)
	if err != nil {
		return err
	}

	if state == InstanceStateRunning {
		return fmt.Errorf("instance %q running already", id)
	}

	runInstance := &ec2.StartInstancesInput{
		InstanceIds: []string{id},
	}
	log.Printf("Start %s", id)
	if outputStart, errInstance := a.Ec2Client.StartInstances(ctx, runInstance); errInstance != nil {
		return err
	} else {
		log.Println(outputStart.StartingInstances)
	}

	return nil
}

func (a AWSCloud) StopInstance(id string) error {
	ctx := context.TODO()
	state, err := a.GetInstanceStatus(id)
	if err != nil {
		return err
	}

	if state != InstanceStateRunning {
		return fmt.Errorf("instance %q not running", id)
	}

	runInstance := &ec2.StopInstancesInput{
		InstanceIds: []string{id},
	}
	log.Printf("Start %s", id)
	if outputStart, errInstance := a.Ec2Client.StopInstances(ctx, runInstance); errInstance != nil {
		return err
	} else {
		log.Println(outputStart.StoppingInstances)
	}

	return nil
}

func (a AWSCloud) GetName() string {
	return "aws"
}

func (a AWSCloud) GetInstanceStatus(id string) (InstanceState, error) {
	ctx := context.TODO()

	var includeAllInstances = true
	input := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: &includeAllInstances,
		InstanceIds:         []string{id},
	}
	output, err := a.Ec2Client.DescribeInstanceStatus(ctx, input)
	if err != nil {
		log.Println(err)
		return "", err
	}

	for _, instanceStatus := range output.InstanceStatuses {
		log.Printf("%s: %s\n", *instanceStatus.InstanceId, instanceStatus.InstanceState.Name)
		if *instanceStatus.InstanceId == id {
			return InstanceState(instanceStatus.InstanceState.Name), nil
		}
	}

	return "", fmt.Errorf("instance status: not found")
}
