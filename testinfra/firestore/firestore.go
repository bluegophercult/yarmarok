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

const (
	defaultPort      = "8080"
	defaultProjectID = "local-test-firestore-project-id"
)

//go:embed Dockerfile
var dockerfile string

// Instance represents a firestore emulator instance
// and controls its lifecycle.
type Instance struct {
	t         *testing.T
	container testcontainers.Container
	ip        string
	client    *firestore.Client
}

// Client returns a firestore client connected to the emulator.
func (i *Instance) Client() *firestore.Client {
	return i.client
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

	t.Setenv("FIRESTORE_EMULATOR_HOST", fmt.Sprintf("%s:%s", ip, "8080"))
	client, err := firestore.NewClient(context.Background(), defaultProjectID)
	if err != nil {
		return nil, fmt.Errorf("create firestore client: %w", err)
	}

	instance := &Instance{
		container: container,
		ip:        ip,
		t:         t,
		client:    client,
	}

	t.Cleanup(func() {
		if err := instance.container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	})

	return instance, err
}