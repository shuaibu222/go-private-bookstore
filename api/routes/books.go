package routes

import (
	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/api/controllers"
	"github.com/shuaibu222/go-bookstore/auth"
)

var BooksRoutes = func(router *mux.Router) {
	router.HandleFunc("/v1/api/books", auth.AuthMiddleware(controllers.GetAllUserBooks)).Methods("GET")
	router.HandleFunc("/v1/api/public/books", controllers.GetAllPublicBooks).Methods("GET")
	router.HandleFunc("/v1/api/books", auth.AuthMiddleware(controllers.CreateNewBook)).Methods("POST")
	router.HandleFunc("/v1/api/books/{id}", auth.AuthMiddleware(controllers.GetBookById)).Methods("GET")
	router.HandleFunc("/v1/api/books/{id}", auth.AuthMiddleware(controllers.UpdateBook)).Methods("PUT")
	router.HandleFunc("/v1/api/books/{id}", auth.AuthMiddleware(controllers.DeleteBook)).Methods("DELETE")
	router.HandleFunc("/v1/api/books/image/{id}", auth.AuthMiddleware(controllers.DeleteCoverFromS3)).Methods("DELETE")
	router.HandleFunc("/v1/api/books/file/{id}", auth.AuthMiddleware(controllers.DeleteBookFromS3)).Methods("DELETE")
}
