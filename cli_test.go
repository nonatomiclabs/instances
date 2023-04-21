package instances_test

import (
	"testing"

	"github.com/nonatomiclabs/instances"
)

func TestCLI(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"no subcommand": {
			args:    []string{},
			wantErr: "use subcommand",
		},
		"unknown subcommand": {
			args:    []string{"johndoe"},
			wantErr: "unknown subcommand",
		},
		"add - no arguments": {
			args:    []string{"add"},
			wantErr: "missing instance ID",
		},
		"add - wrong arguments": {
			args:    []string{"add", "--cloud-provider", "aws"},
			wantErr: "flag provided but not defined",
		},
		"add - existing instance name": {
			args:    []string{"add", "--name", existingInstanceName, "--cloud", "mock", existingInstanceIds[1]},
			wantErr: "exists already",
		},
		"add - existing instance id": {
			args:    []string{"add", "--name", "testInstance", "--cloud", "mock", existingInstanceIds[0]},
			wantErr: "already referenced",
		},
		"add - multiple instance ids": {
			args:    []string{"add", "--name", "testInstance", "--cloud", "mock", "anInstanceId", "anotherInstanceId"},
			wantErr: "only one instance",
		},
		"add - nonexisting cloud provider": {
			args:    []string{"add", "--name", "testInstance", "--cloud", "myGreatCloud", existingInstanceIds[1]},
			wantErr: "unsupported cloud provider",
		},
		"add - new instance": {
			args:    []string{"add", "--name", "testInstance", "--cloud", "mock", existingInstanceIds[1]},
			wantErr: "",
		},
		"remove - existing instance": {
			args:    []string{"rm", existingInstanceName},
			wantErr: "",
		},
		"remove - nonexisting instance": {
			args:    []string{"rm", "anInstance"},
			wantErr: "no instance named",
		},
		"remove - no arguments": {
			args:    []string{"rm"},
			wantErr: "missing instance name",
		},
		"remove - multiple instances": {
			args:    []string{"rm", "anInstance", "anotherInstance"},
			wantErr: "only one instance",
		},
		"remove - wrong flags": {
			args:    []string{"rm", "--option", "value"},
			wantErr: "flag provided but not defined",
		},
		"status - existing instance": {
			args:    []string{"status", existingInstanceName},
			wantErr: "",
		},
		"status - nonexisting instance": {
			args:    []string{"status", "anInstance"},
			wantErr: "no instance named",
		},
		"status - no arguments": {
			args:    []string{"status"},
			wantErr: "missing instance name",
		},
		"status - wrong arguments": {
			args:    []string{"status", "--option", "value"},
			wantErr: "flag provided but not defined",
		},
		"start - existing instance": {
			args:    []string{"start", existingInstanceName},
			wantErr: "",
		},
		"start - nonexisting instance": {
			args:    []string{"start", "anInstance"},
			wantErr: "no instance named",
		},
		"start - no arguments": {
			args:    []string{"start"},
			wantErr: "missing instance name",
		},
		"start - wrong arguments": {
			args:    []string{"start", "--option", "value"},
			wantErr: "flag provided but not defined",
		},
		"stop - existing instance": {
			args:    []string{"stop", existingInstanceName},
			wantErr: "",
		},
		"stop - nonexisting instance": {
			args:    []string{"stop", "anInstance"},
			wantErr: "no instance named",
		},
		"stop - no arguments": {
			args:    []string{"stop"},
			wantErr: "missing instance name",
		},
		"stop - wrong arguments": {
			args:    []string{"stop", "--option", "value"},
			wantErr: "flag provided but not defined",
		},
		"list - no arguments": {
			args:    []string{"list"},
			wantErr: "",
		},
		"list - arguments": {
			args:    []string{"list", existingInstanceName},
			wantErr: "list doesn't take positional arguments",
		},
		"list - wrong options": {
			args:    []string{"list", "--option", "value"},
			wantErr: "flag provided but not defined",
		},
		"list - valid options": {
			args:    []string{"list", "--cloud", "mock"},
			wantErr: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := getInitializedDatabase()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			cloudProviders := map[string]instances.CloudProvider{
				"mock": MockCloudProvider{},
			}

			cli := instances.NewCLI(db, cloudProviders)

			err = cli.Run(test.args)
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
