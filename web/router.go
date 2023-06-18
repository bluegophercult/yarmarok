package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kaznasho/yarmarok/logger"
	"github.com/kaznasho/yarmarok/service"
)

const (
	RafflesPath      = "/raffles"
	ContributorsPath = "/contributors"

	raffleIDParam      = "raffle_id"
	contributorIDParam = "contributor_id"

	raffleIDPlaceholder = "/{" + raffleIDParam + "}"
)

var (
	// ErrAmbiguousOrganizerIDHeader is returned when
	// the organizer id header is not set or is ambiguous.
	ErrAmbiguousOrganizerIDHeader = errors.New("ambiguous organizer id format")

	// ErrMissingID is returned when id is missing.
	ErrMissingID = errors.New("missing id")
)

// Router is responsible for routing requests
// to the corresponding services.
type Router struct {
	chi.Router
	organizerService service.OrganizerService
	logger           *logger.Entry
}

// NewRouter creates a new Router
func NewRouter(us service.OrganizerService, log *logger.Logger) (*Router, error) {
	router := &Router{
		Router:           chi.NewRouter(),
		organizerService: us,
		logger: log.WithFields(
			logger.Fields{
				"component": "router",
				"trace_id":  uuid.New().String(),
			},
		),
	}

	router.Use(router.corsMiddleware)
	router.Use(router.loggingMiddleware)
	router.Use(router.recoverMiddleware)
	router.Use(router.organizerMiddleware)

	// "/raffles"
	router.Route(RafflesPath, func(rafRouter chi.Router) {
		rafRouter.Post("/", router.createRaffle)
		rafRouter.Get("/", router.listRaffles)
		// "/raffles/{raffle_id}"
		rafRouter.Route(raffleIDPlaceholder, func(rafIDRouter chi.Router) {
			// "/raffles/{raffle_id}/contributors"
			rafIDRouter.Route(ContributorsPath, func(ctbRouter chi.Router) {
				ctbRouter.Post("/", router.createContributor)
				ctbRouter.Put("/", router.updateContributor)
				ctbRouter.Get("/", router.listContributors)
			})
		})
	})

	return router, nil
}

func (r *Router) createRaffle(w http.ResponseWriter, req *http.Request) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	raffleService := r.organizerService.RaffleService(organizerID)

	m := newMethodHandler(raffleService.Init, r.logger.Logger)

	m.ServeHTTP(w, req)
}

func (r *Router) listRaffles(w http.ResponseWriter, req *http.Request) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	raffleService := r.organizerService.RaffleService(organizerID)

	m := newNoRequestMethodHandler(raffleService.List, r.logger.Logger)

	m.ServeHTTP(w, req)
}

func (r *Router) createContributor(w http.ResponseWriter, req *http.Request) {
	contributorService, err := r.getContributorService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMethodHandler(contributorService.Add, r.logger.Logger).ServeHTTP(w, req)
}

func (r *Router) updateContributor(w http.ResponseWriter, req *http.Request) {
	contributorService, err := r.getContributorService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMethodHandler(contributorService.Edit, r.logger.Logger).ServeHTTP(w, req)
}

func (r *Router) listContributors(w http.ResponseWriter, req *http.Request) {
	contributorService, err := r.getContributorService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newNoRequestMethodHandler(contributorService.List, r.logger.Logger).ServeHTTP(w, req)
}

func (r *Router) getContributorService(req *http.Request) (service.ContributorService, error) {
	organizerID, err := extractOrganizerID(req)
	if err != nil {
		return nil, err
	}

	raffleID, err := extractParam(req, raffleIDParam)
	if err != nil || raffleID == "" {
		return nil, ErrMissingID
	}

	contributorService := r.organizerService.RaffleService(organizerID).ContributorService(raffleID)

	return contributorService, nil
}

func extractOrganizerID(r *http.Request) (string, error) {
	ids := r.Header.Values(GoogleOrganizerIDHeader)

	if len(ids) != 1 {
		return "", ErrAmbiguousOrganizerIDHeader
	}

	id := ids[0]
	if id == "" {
		return "", ErrAmbiguousOrganizerIDHeader
	}

	return id, nil
}

func extractParam(req *http.Request, param string) (string, error) {
	val := chi.URLParam(req, param)
	if param == "" {
		return "", fmt.Errorf("missing param: %s", param)
	}

	return val, nil
}
