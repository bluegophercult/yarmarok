package function

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
	"github.com/kaznasho/yarmarok/storage"
	"github.com/kaznasho/yarmarok/web"
)

const projectIDEnvVar = "GCP_PROJECT"

// ErrEmptyProjectID is returned when the project id is empty.
var ErrEmptyProjectID = errors.New("empty project id")

// EntryPoint is the entry point for the cloud function.
func Entrypoint(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger(logger.LevelInfo)

	router, err := LoadRouter(log)
	if err != nil {
		log.WithField("component", "entrypoint").Error("Loading router error: ", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	router.ServeHTTP(w, r)
}

// LoadRouter loads the router.
func LoadRouter(log *logger.Logger) (*web.Router, error) {
	projectID := os.Getenv(projectIDEnvVar)
	if projectID == "" {
		return nil, fmt.Errorf("%w: %s is not set", ErrEmptyProjectID, projectIDEnvVar)
	}

	firestoreClient, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}

	userStorage := storage.NewFirestoreUserStorage(firestoreClient)

	userService := service.NewUserManager(userStorage)

	return web.NewRouter(userService, log)
}
