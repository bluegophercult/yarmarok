package testinfra

import (
	"os"
	"testing"
)

func SkipIfNotIntegrationRun(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}
}
