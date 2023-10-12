package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/api/routes"
)

func main() {
	r := mux.NewRouter()

	routes.BooksRoutes(r)
	routes.UsersRoutes(r)
	http.Handle("/", r)
	fmt.Println("Server started.........")

	port := os.Getenv("WEB_PORT")

	log.Fatal(http.ListenAndServe("localhost:"+port,
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedOrigins([]string{"http://*", "https://*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"*"}),
		)(r)),
	)
}
