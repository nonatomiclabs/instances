package instances_test

import (
	"bytes"
	"errors"
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

var existingInstanceIds = []string{"existingInstance1", "existingInstance2"}

const existingInstanceName = "myInstance"

type MockCloudProvider struct{}

func (m MockCloudProvider) GetInstanceStatus(id string) (instances.InstanceState, error) {
	switch id {
	case existingInstanceIds[0], existingInstanceIds[1]:
		return instances.InstanceStateRunning, nil
	default:
		return "", fmt.Errorf("instance %q not found in the cloud provider", id)
	}
}

func (m MockCloudProvider) StartInstance(id string) error {
	return nil
}

func (m MockCloudProvider) StopInstance(id string) error {
	return nil
}

func (m MockCloudProvider) GetName() string {
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
			instanceId: existingInstanceIds[0],
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

			err = db.AddInstance(test.instanceId, "instanceName", MockCloudProvider{})

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
		"instance does not exist yet": {
			instanceName: "aNewInstance",
			wantErr:      "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.NewBufferString("{\"instances\": {}}")
			db, err := instances.NewDatabase(buf)
			if err != nil {
				t.Fatalf("could not acquire db: %s", err)
			}

			err = db.AddInstance(existingInstanceIds[0], "alreadyPresent", MockCloudProvider{})
			if err != nil {
				t.Fatal("failed to add pre-required instance")
			}

			err = db.AddInstance(existingInstanceIds[1], test.instanceName, MockCloudProvider{})
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAddInstanceIdCheck(t *testing.T) {
	tests := map[string]struct {
		instanceId string
		wantErr    string
	}{
		"instance exists already in the database": {
			instanceId: existingInstanceIds[0],
			wantErr:    "already referenced",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := getInitializedDatabase()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = db.AddInstance(test.instanceId, "testInstance", MockCloudProvider{})
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func getInitializedDatabase() (*instances.Database, error) {
	buf := bytes.NewBufferString("{\"instances\": {}}")
	db, err := instances.NewDatabase(buf)
	if err != nil {
		return nil, fmt.Errorf("could not acquire db: %s", err)
	}

	err = db.AddInstance(existingInstanceIds[0], existingInstanceName, MockCloudProvider{})
	if err != nil {
		return nil, errors.New("failed to add pre-required instance")
	}
	return db, nil
}

func TestGetInstance(t *testing.T) {
	tests := map[string]struct {
		instanceName string
		wantErr      string
	}{
		"existing instance": {
			instanceName: existingInstanceName,
			wantErr:      "",
		},
		"nonexisting instance": {
			instanceName: "iDontExist",
			wantErr:      "no instance named",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := getInitializedDatabase()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			_, err = db.GetInstance(test.instanceName)
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRemoveInstance(t *testing.T) {
	tests := map[string]struct {
		instanceName string
		wantErr      string
	}{
		"existing instance": {
			instanceName: existingInstanceName,
			wantErr:      "",
		},
		"nonexisting instance": {
			instanceName: "iDontExist",
			wantErr:      "no instance named",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := getInitializedDatabase()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = db.RemoveInstance(test.instanceName)
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSaveDatabase(t *testing.T) {
	db, err := getInitializedDatabase()
	if err != nil {
		t.Fatalf("test setup failed: %v", err)
	}

	var bufOut bytes.Buffer
	err = db.Save(&bufOut)
	if err != nil {
		t.Fatal(err)
	}

	if bufOut.Len() == 0 {
		t.Fatal("saving database did not write any content")
	}
}
