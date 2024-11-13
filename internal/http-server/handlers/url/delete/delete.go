package delete

import (
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/response"
	sl "url-shortener/internal/lib/logger/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// Request contains the alias of the URL to be deleted
type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type URLDelete interface {
	DeleteURL(alias string) error
}

// New creates a new HTTP handler function for deleting URLs by alias.
// It extracts the alias from the request path and calls the provided URLDelete
// interface to delete the URL. If the alias is empty or the deletion fails, it
// responds with an appropriate JSON response. Otherwise, it responds with a 200 OK response.
func New(log *slog.Logger, delete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.url.delete.New"
		log = log.With(slog.String("operation", operation))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			slog.Info("missing alias in request")
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		log.Info("recieved request to delete url", slog.String("alias", alias))

		if err := delete.DeleteURL(alias); err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}

		render.JSON(w, r, resp.OK())
	}
}

// validateRequest validates the request
func validateRequest(req Request) error {
	validate := validator.New()
	return validate.Struct(req)
}
