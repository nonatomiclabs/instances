package instances_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/nonatomiclabs/instances"
)

// ErrorContains checks if the error message in out contains the text in
// want.
func errorContains(err error, want string) bool {
	if err == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(err.Error(), want)
}

const existingInstanceId = "exists"

type mockCloudProvider struct{}

func (m mockCloudProvider) GetInstanceStatus(id string) (instances.InstanceState, error) {
	switch id {
	case existingInstanceId:
		return instances.InstanceStateRunning, nil
	default:
		return "", fmt.Errorf("instance %q not found in the cloud provider", id)
	}
}

func (m mockCloudProvider) StartInstance(id string) error {
	return nil
}

func (m mockCloudProvider) StopInstance(id string) error {
	return nil
}

func (m mockCloudProvider) GetName() string {
	return ""
}

func TestNewDatabase(t *testing.T) {
	tests := map[string]struct {
		payload string
		wantErr string
	}{
		"valid payload": {
			payload: "{\"instances\": {}}",
			wantErr: "",
		},
		"empty payload": {
			payload: "",
			wantErr: "EOF",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBufferString(test.payload)
			_, err := instances.NewDatabase(buf)

			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAddInstanceCloudProviderCheck(t *testing.T) {
	tests := map[string]struct {
		instanceId string
		wantErr    string
	}{
		"instance exists in cloud provider": {
			instanceId: existingInstanceId,
			wantErr:    "",
		},
		"instance does not exist in cloud provider": {
			instanceId: "noExists",
			wantErr:    "not found in the cloud provider",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBufferString("{\"instances\": {}}")
			db, err := instances.NewDatabase(buf)
			if err != nil {
				t.Fatalf("could not acquire db: %s", err)
			}

			err = db.AddInstance(test.instanceId, "instanceName", mockCloudProvider{})

			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAddInstanceNameCheck(t *testing.T) {
	tests := map[string]struct {
		instanceName string
		wantErr      string
	}{
		"instance exists already in the database": {
			instanceName: "alreadyPresent",
			wantErr:      "exists already",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBufferString("{\"instances\": {}}")
			db, err := instances.NewDatabase(buf)
			if err != nil {
				t.Fatalf("could not acquire db: %s", err)
			}

			err = db.AddInstance(existingInstanceId, "alreadyPresent", mockCloudProvider{})
			if err != nil {
				t.Fatal("failed to add pre-required instance")
			}

			err = db.AddInstance(existingInstanceId, test.instanceName, mockCloudProvider{})
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
