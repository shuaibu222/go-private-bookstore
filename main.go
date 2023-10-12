package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/api/routes"
	config "github.com/shuaibu222/go-bookstore/config"
)

func main() {
	r := mux.NewRouter()

	routes.BooksRoutes(r)
	routes.UsersRoutes(r)
	http.Handle("/", r)
	fmt.Println("Server started.........")

	config, err := config.LoadConfig()
	if err != nil {
		log.Println("Error while loading envs: ", err)
	}

	log.Fatal(http.ListenAndServe("localhost:"+config.WebPort,
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedOrigins([]string{"http://*", "https://*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"*"}),
		)(r)),
	)
}
