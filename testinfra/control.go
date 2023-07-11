package testinfra

import (
	"os"
	"strconv"
	"testing"
)

// SkipIfNotIntegrationRun skips the test if RUN_INTEGRATION_TESTS is not set to true.
func SkipIfNotIntegrationRun(t *testing.T) {
	if ok, _ := strconv.ParseBool(os.Getenv("RUN_INTEGRATION_TESTS")); !ok {
		t.Skip("Skipping integration test")
	}
}
