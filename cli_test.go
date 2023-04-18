package instances_test

import (
	"testing"

	"github.com/nonatomiclabs/instances"
)

func TestAddInstance(t *testing.T) {
	db, err := getInitializedDatabase()
	if err != nil {
		t.Fatalf("test setup failed: %v", err)
	}

	cloudProviders := map[string]instances.CloudProvider{
		"mock": MockCloudProvider{},
	}

	cli := instances.NewCLI(db, cloudProviders)

	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"no arguments": {
			args:    []string{},
			wantErr: "missing instance ID",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err = cli.Run(append([]string{"add"}, test.args...))
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
