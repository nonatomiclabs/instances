package instances_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nonatomiclabs/instances"
)

const runningInstanceId = "i-1234"
const nonRunningInstanceId = "i-5678"

type mockEC2Manager struct{}

func (m mockEC2Manager) DescribeInstanceStatus(ctx context.Context, params *ec2.DescribeInstanceStatusInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error) {
	out := ec2.DescribeInstanceStatusOutput{}

	if len(params.InstanceIds) == 1 {
		id := params.InstanceIds[0]
		switch id {
		case runningInstanceId:
			status := types.InstanceStatus{
				InstanceState: &types.InstanceState{Name: types.InstanceStateNameRunning},
				InstanceId:    &id,
			}
			out.InstanceStatuses = append(out.InstanceStatuses, status)
		}
	}

	return &out, nil
}

func (m mockEC2Manager) StartInstances(ctx context.Context, params *ec2.StartInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error) {
	return &ec2.StartInstancesOutput{}, nil
}

func TestStartEC2Instance(t *testing.T) {
	tests := map[string]struct {
		instanceID string
		wantErr    string
	}{
		"running instance": {
			instanceID: runningInstanceId,
			wantErr:    "running already",
		},
		"non-running instance": {
			instanceID: nonRunningInstanceId,
			wantErr:    "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := mockEC2Manager{}
			AWSCloud := instances.AWSCloud{Ec2Client: client}
			err := AWSCloud.StartInstance(test.instanceID)
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}

}
