package main

import (
	"fmt"
	"net/http"

	"github.com/000xs/htmx-todo-app/routes"
	"github.com/gorilla/handlers"
)

const (
	VERSION = "1.0.0"
	PORT    = "3000"
)

func main() {
	r := routes.NewRouter()

	go fmt.Printf("%vServer is Running on http://localhost:%v%v   \n", Magenta, PORT, Magenta)
	http.ListenAndServe(":3000", handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)(r))

}
