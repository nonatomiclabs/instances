package instances_test

import (
	"testing"

	"github.com/nonatomiclabs/instances"
)

func getInitializedCLI() (*instances.CLI, error) {
	db, err := getInitializedDatabase()
	if err != nil {
		return nil, err
	}

	cloudProviders := map[string]instances.CloudProvider{
		"mock": MockCloudProvider{},
	}

	return instances.NewCLI(db, cloudProviders), nil
}

func TestCLI(t *testing.T) {
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cli, err := getInitializedCLI()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = cli.Run(test.args)
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAddInstanceCLI(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"no arguments": {
			args:    []string{},
			wantErr: "missing instance ID",
		},
		"wrong arguments": {
			args:    []string{"--cloud-provider", "aws"},
			wantErr: "flag provided but not defined",
		},
		"existing instance name": {
			args:    []string{"--name", existingInstanceName, "--cloud", "mock", existingInstanceIds[1]},
			wantErr: "exists already",
		},
		"existing instance id": {
			args:    []string{"--name", "testInstance", "--cloud", "mock", existingInstanceIds[0]},
			wantErr: "already referenced",
		},
		"multiple instance ids": {
			args:    []string{"--name", "testInstance", "--cloud", "mock", "anInstanceId", "anotherInstanceId"},
			wantErr: "only one instance",
		},
		"nonexisting cloud provider": {
			args:    []string{"--name", "testInstance", "--cloud", "myGreatCloud", existingInstanceIds[1]},
			wantErr: "unsupported cloud provider",
		},
		"successful add": {
			args:    []string{"--name", "testInstance", "--cloud", "mock", existingInstanceIds[1]},
			wantErr: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cli, err := getInitializedCLI()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = cli.Run(append([]string{"add"}, test.args...))
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRemoveInstanceCLI(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"existing instance": {
			args:    []string{existingInstanceName},
			wantErr: "",
		},
		"nonexisting instance": {
			args:    []string{"anInstance"},
			wantErr: "no instance named",
		},
		"no arguments": {
			args:    []string{},
			wantErr: "missing instance name",
		},
		"multiple instances": {
			args:    []string{"anInstance", "anotherInstance"},
			wantErr: "only one instance",
		},
		"wrong flags": {
			args:    []string{"--option", "value"},
			wantErr: "flag provided but not defined",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cli, err := getInitializedCLI()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = cli.Run(append([]string{"rm"}, test.args...))
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}

		})
	}
}

func TestGetInstanceStatusCLI(t *testing.T) {
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"valid instance": {
			args:    []string{existingInstanceName},
			wantErr: "",
		},
		"nonexisting instance": {
			args:    []string{"anInstance"},
			wantErr: "no instance named",
		},
		"no arguments": {
			args:    []string{},
			wantErr: "missing instance name",
		},
		"wrong arguments": {
			args:    []string{"--option", "value"},
			wantErr: "flag provided but not defined",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cli, err := getInitializedCLI()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = cli.Run(append([]string{"status"}, test.args...))
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}

		})
	}
}

func TestStartInstance(t *testing.T) {
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"valid instance": {
			args:    []string{existingInstanceName},
			wantErr: "",
		},
		"invalid instance": {
			args:    []string{"anInstance"},
			wantErr: "no instance named",
		},
		"no arguments": {
			args:    []string{},
			wantErr: "missing instance name",
		},
		"wrong arguments": {
			args:    []string{"--option", "value"},
			wantErr: "flag provided but not defined",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cli, err := getInitializedCLI()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = cli.Run(append([]string{"start"}, test.args...))
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestStopInstance(t *testing.T) {
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"valid instance": {
			args:    []string{existingInstanceName},
			wantErr: "",
		},
		"invalid instance": {
			args:    []string{"anInstance"},
			wantErr: "no instance named",
		},
		"no arguments": {
			args:    []string{},
			wantErr: "missing instance name",
		},
		"wrong arguments": {
			args:    []string{"--option", "value"},
			wantErr: "flag provided but not defined",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cli, err := getInitializedCLI()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = cli.Run(append([]string{"stop"}, test.args...))
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestListInstance(t *testing.T) {
	tests := map[string]struct {
		args    []string
		wantErr string
	}{
		"no arguments": {
			args:    []string{},
			wantErr: "",
		},
		"wrong arguments": {
			args:    []string{"--option", "value"},
			wantErr: "flag provided but not defined",
		},
		"valid arguments": {
			args:    []string{"--cloud", "mock"},
			wantErr: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cli, err := getInitializedCLI()
			if err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			err = cli.Run(append([]string{"list"}, test.args...))
			if !errorContains(err, test.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
