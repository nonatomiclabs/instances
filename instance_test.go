package instances_test

import (
	"testing"

	"github.com/nonatomiclabs/instances"
)

func TestGetCloudProvider(t *testing.T) {
	cloudProviders := map[string]instances.CloudProvider{
		"mock": MockCloudProvider{},
	}

	tests := map[string]struct {
		instanceCloudProvider string
		want                  instances.CloudProvider
		wantErr               string
	}{
		"existing cloud provider": {
			instanceCloudProvider: "mock",
			want:                  MockCloudProvider{},
			wantErr:               "",
		},
		"nonexisting cloud provider": {
			instanceCloudProvider: "aws",
			want:                  nil,
			wantErr:               "unsupported cloud provider",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			instance := instances.Instance{
				Id:                "abcd",
				CloudProviderName: test.instanceCloudProvider,
			}

			cloudProvider, err := instance.GetCloudProvider(cloudProviders)
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}

			if cloudProvider != test.want {
				t.Fatalf("wrong cloud provider: got %v, want %v", cloudProvider, test.want)
			}

		})
	}

}
