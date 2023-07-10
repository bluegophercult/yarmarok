//go:build local

package local

import (
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/stretchr/testify/require"

	_ "github.com/kaznasho/yarmarok"
	"github.com/kaznasho/yarmarok/function"
	"github.com/kaznasho/yarmarok/testinfra/firestore"
)

// TestRun used to run service locally
func TestRun(t *testing.T) {
	firestoreInstance, err := firestore.RunInstance(t)
	require.NoError(t, err)

	t.Setenv(function.ProjectIDEnvVar, firestoreInstance.ProjectID())
	t.Setenv("FUNCTION_TARGET", "Entrypoint")

	port := "8081"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	t.Logf("Starting on port %s\n", port)
	if err = funcframework.Start(port); err != nil {
		t.Fatalf("Start func: %s\n", err)
	}
}
