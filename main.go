package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/api/routes"
)

func main() {
	r := mux.NewRouter()

	routes.BooksRoutes(r)
	routes.UsersRoutes(r)
	http.Handle("/", r)
	fmt.Println("Server started at http://localhost:9000 ......")

	log.Fatal(http.ListenAndServe("localhost:9000",
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedOrigins([]string{"http://*", "https://*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"*"}),
		)(r)),
	)
}
