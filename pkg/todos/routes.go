package todos

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Handler(svc Service) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Use(middleware.Timeout(time.Minute))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

	r.With(middleware.AllowContentType("application/json")).Route("/todos", func(r chi.Router) {
		r.Get("/", listTodos(svc))
		r.Post("/", createTodo(svc))

		r.Route("/{id:[0-9a-zA-Z-]+}", func(r chi.Router) {
			r.Use(TodoCtx(svc))
			r.Get("/", getTodo(svc))
			r.Put("/", updateTodo(svc))
			r.Delete("/", deleteTodo(svc))
		})
	})

	return r
}

type todoCtxKey struct{}

var TodoCtxKey = &todoCtxKey{}

func TodoCtx(svc Service) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			id := chi.URLParam(r, "id")
			td, err := svc.FindByID(r.Context(), id)
			if err != nil {
				log.Printf("Preparing TodoCtx: %v\n", err)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), TodoCtxKey, &td)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
