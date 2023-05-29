// package firestore provides a firestore emulator instance for testing.
package firestore

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

//go:embed Dockerfile
var dockerfile string

// Instance represents a firestore emulator instance
// and controls its lifecycle.
type Instance struct {
	t         *testing.T
	container testcontainers.Container
	ip        string
}

// IP returns the ip of the emulator container.
func (i *Instance) IP() string {
	return i.ip
}

// Port returns the port to connect to.
func (i *Instance) Port() string {
	return "8080"
}

// Address returns the address of the emulator container.
func (i *Instance) Address() string {
	return fmt.Sprintf("%s:%s", i.IP(), i.Port())
}

// ProjectID returns the project id of the emulator.
func (i *Instance) ProjectID() string {
	return "local-test-firestore-project-id"
}

// Client returns a firestore client connected to the emulator.
func (i *Instance) Client() (*firestore.Client, error) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", i.Address())
	return firestore.NewClient(context.Background(), i.ProjectID())
}

// RunInstance runs a firestore emulator instance in a docker container.
// The container is automatically cleaned up after the test.
// Use only one instance per host, as the emulator does not support multiple instances.
func RunInstance(t *testing.T) (*Instance, error) {
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
	//waitStrategy := (&wait.NopStrategy{}).WithStartupTimeout(5 * time.Second)
	req := testcontainers.ContainerRequest{
		FromDockerfile: fromDockerfile,
		ExposedPorts:   []string{"8080/tcp"},
		WaitingFor:     wait.ForExposedPort(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}

	ip, err := container.ContainerIP(ctx)
	if err != nil {
		return nil, fmt.Errorf("get container host: %w", err)
	}

	instance := &Instance{
		container: container,
		ip:        ip,
		t:         t,
	}

	t.Cleanup(func() {
		if err := instance.container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	})

	return instance, err
}
