// Package firestore provides a firestore emulator instance for testing.
package firestore

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultPort      = "8080"
	defaultProjectID = "local-test-firestore-project-id"
)

//go:embed Dockerfile
var dockerfile string

// Instance represents a firestore emulator instance
// and controls its lifecycle.
type Instance struct {
	client    *firestore.Client
	container testcontainers.Container
}

// Client returns a firestore client connected to the emulator.
func (i *Instance) Client() *firestore.Client {
	return i.client
}

// ProjectID returns the project id of the emulator.
func (i *Instance) ProjectID() string {
	return defaultProjectID
}

type TestEnv interface {
	TempDir() string
	Cleanup(f func())
	Fatalf(format string, args ...any)
	Log(args ...any)
	Setenv(key, value string)
}

// RunInstance runs a firestore emulator instance in a docker container.
// The container is automatically cleaned up after the test.
// Use only one instance per host, as the emulator does not support multiple instances.
func RunInstance(t TestEnv) (*Instance, error) {
	tempDir := t.TempDir()
	dockerDir := path.Join(tempDir, uuid.New().String())

	err := os.MkdirAll(dockerDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("create docker dir: %w", err)
	}

	err = os.WriteFile(path.Join(dockerDir, "Dockerfile"), []byte(dockerfile), 0644)
	if err != nil {
		return nil, fmt.Errorf("copy dockerfile: %w", err)
	}

	fromDockerfile := testcontainers.FromDockerfile{
		Context: dockerDir,
	}

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Name:           "yarmarok-firestore",
		FromDockerfile: fromDockerfile,
		ExposedPorts:   []string{fmt.Sprintf("%s:%s", defaultPort, "8080/tcp")},
		WaitingFor:     wait.ForExposedPort(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}

	t.Setenv("FIRESTORE_EMULATOR_HOST", fmt.Sprintf("%s:%s", "localhost", defaultPort))
	client, err := firestore.NewClient(context.Background(), defaultProjectID)
	if err != nil {
		return nil, fmt.Errorf("create firestore client: %w", err)
	}

	instance := &Instance{
		container: container,
		client:    client,
	}

	t.Cleanup(func() {
		if err = instance.container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate firestore container: %s", err)
		}
		t.Log("Stopped firestore container")
	})

	return instance, err
}
