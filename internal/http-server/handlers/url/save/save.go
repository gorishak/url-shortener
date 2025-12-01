package save

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/gorishak/url-shortener/internal/lib/api/response"
	"github.com/gorishak/url-shortener/internal/lib/logger/sl"
	"github.com/gorishak/url-shortener/internal/lib/random"
	"github.com/gorishak/url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	validator "github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const (
	aliasLength = 6
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLSaver --output=./mocks --outpkg=mocks --filename=urlsaver_mock.go
type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		// an := make(map[string]interface{})
		// read, _ := io.ReadAll(r.Body)
		// _ = json.Unmarshal(read, &an)
		// fmt.Println(an)
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("validation failed", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url saved", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
