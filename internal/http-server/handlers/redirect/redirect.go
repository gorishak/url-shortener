package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/gorishak/url-shortener/internal/lib/api/response"
	"github.com/gorishak/url-shortener/internal/lib/logger/sl"
	"github.com/gorishak/url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// url getter interface for getting url by alias

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLGetter --output=./mocks --outpkg=mocks --filename=urlgetter_mock.go
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resUrl, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}
		log.Info("url got found", slog.String("url", resUrl))

		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
