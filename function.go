package yarmarok

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/kaznasho/yarmarok/function"
)

func init() {
	// Register an HTTP function with the Functions Framework
	functions.HTTP("Entrypoint", function.Entrypoint)
}
