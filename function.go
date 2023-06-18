package yarmarok

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/kaznasho/yarmarok/function"
)

func init() {
	functions.HTTP("Entrypoint", Entrypoint)
}

var Entrypoint = function.Entrypoint
