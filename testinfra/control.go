package testinfra

import (
	"os"
	"testing"
)

// SkipIfNotIntegrationRun skips the test if RUN_INTEGRATION_TESTS is not set to true.
func SkipIfNotIntegrationRun(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test")
	}
}
