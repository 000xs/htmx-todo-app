package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// NewRouter initializes and returns a new router
func NewRouter() *chi.Mux {

	// Wrap router with CORS

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/", RootHandler)
	r.Get("/register", RegisterFrontHandler)
	r.Get("/login", LoginFrontHandler)

	r.Get("/api/todo/{userId}", TodoGetHandler)
	r.Post("/api/todo", TodoPostHandler)
	r.Put("/api/todo", TodoUpdateHandler)

	//auth
	r.Post("/api/auth/register", RegisterHandler)
	r.Post("/api/auth/login", LoginHandler)
	r.Get("/api/auth/validate", ValidateHandler)
	return r
}

func reqLog(r *http.Request) {

	fmt.Printf("%v %v  -  %v - %v %v %v\n", "\033[36m", r.Method, r.URL, r.Proto, r.UserAgent(), "\033[36m")

}
