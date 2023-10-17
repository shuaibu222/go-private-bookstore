package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/models"
	"github.com/shuaibu222/go-bookstore/utils"
)

var mutex sync.Mutex
var uploadWg sync.WaitGroup
var paths []string

func CreateNewBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	id, userName := utils.JwtUserIdUsername(w, r)

	CreateBook := &models.Books{}

	// anonymous field to insert user's Id immediately before creating a book instance
	CreateBook.User.UserId = id
	CreateBook.User.Username = userName

	// Parse the file from the request
	book, book_header, err := r.FormFile("uploaded_book")
	if err != nil {
		log.Println("Error uploading book: ", err)
		http.Error(w, "Failed to read book from request", http.StatusBadRequest)
		return
	}

	cover, cover_header, err := r.FormFile("uploaded_cover_image")
	if err != nil {
		log.Println("Error uploading cover image: ", err)
		http.Error(w, "Failed to read cover image from request", http.StatusBadRequest)
		return
	}

	defer book.Close()
	defer cover.Close()

	// parse the book instance
	json.NewDecoder(r.Body).Decode(&CreateBook)

	// Clear the paths slice to avoid attaching previous URLs
	paths = nil

	books := models.GetAllBooks(id)

	for _, book := range books {
		if CreateBook.Title == book.Title && CreateBook.AuthorName == book.AuthorName {
			json.NewEncoder(w).Encode("This book already exists. No duplicate books")
			return
		}
	}

	// Creating a channel to receive the results
	resultChan := make(chan string, 2) // Buffer size 2 for two goroutines

	uploadWg.Add(2)

	// is better to call uploadWg.Done() in our goroutine initiator function

	go func() {
		defer uploadWg.Done()
		utils.UploadBook(book_header.Filename, resultChan)
	}()

	go func() {
		defer uploadWg.Done()
		// You can at least upload a book even if you don't upload a book cover
		if path.Ext(CreateBook.UploadedBook) != ".pdf" && path.Ext(CreateBook.UploadedBook) != ".txt" {
			json.NewEncoder(w).Encode("You must at least upload a book file")
			return
		} else if path.Ext(CreateBook.UploadedBook) == ".pdf" || path.Ext(CreateBook.UploadedBook) == ".txt" {
			utils.UploadCoverImage(cover_header.Filename, resultChan)
		}
	}()

	uploadWg.Wait()

	// Close the result channel since we're done with it
	close(resultChan)

	// append all end result in paths variable
	for result := range resultChan {
		mutex.Lock()
		paths = append(paths, result)
		mutex.Unlock()
	}

	for _, pathlink := range paths {
		if path.Ext(pathlink) == ".pdf" {
			CreateBook.BookURL.FileURL = pathlink
		} else if path.Ext(pathlink) == ".txt" {
			CreateBook.BookURL.FileURL = pathlink
		} else {
			CreateBook.BookURL.CoverImageUrl = pathlink
		}
	}
	if CreateBook.FileURL == "" {
		json.NewEncoder(w).Encode("Failed to create book!")
	} else {
		// create a new book instance
		book := CreateBook.CreateBook()
		res, _ := json.Marshal(book)
		w.Write(res)

	}

}

func GetAllUserBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	id, _ := utils.JwtUserIdUsername(w, r)

	books := models.GetAllBooks(id)

	res, _ := json.Marshal(books)
	w.Write(res)
}

func GetAllPublicBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	books := models.GetPublicBooks()

	res, _ := json.Marshal(books)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)
	ID, err := strconv.ParseInt(params["id"], 0, 0) // convert to int
	if err != nil {
		log.Println("Error while parsing!")
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	founded, _ := models.GetBookById(ID)

	if founded.Privacy && founded.UserId == id {
		res, _ := json.Marshal(founded)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to view this book!")
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var updateBook = &models.Books{} // initialize a empty book struct to hold the updated values
	utils.ParseBody(r, updateBook)   // parse the body taken from client request for golang to understand and forward it to the DB
	params := mux.Vars(r)
	bookId := params["id"]
	ID, err := strconv.ParseInt(bookId, 0, 0) // convert to int
	if err != nil {
		log.Println("Error while parsing!")
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	bookDetails, db := models.GetBookById(ID)

	if bookDetails.UserId == id {
		if updateBook.Title != "" {
			bookDetails.Title = updateBook.Title
		}
		if updateBook.Description != "" {
			bookDetails.Description = updateBook.Description
		}

		if updateBook.AuthorName != "" {
			bookDetails.AuthorName = updateBook.AuthorName
		}
		if updateBook.AuthorBio != "" {
			bookDetails.AuthorBio = updateBook.AuthorBio
		}
		if updateBook.PublishDate != "" {
			bookDetails.PublishDate = updateBook.PublishDate
		}
		if updateBook.Genre != "" {
			bookDetails.Genre = updateBook.Genre
		}
		if updateBook.UploadedBook != "" {
			bookDetails.Genre = updateBook.UploadedBook
		}
		if updateBook.UploadedCoverImage != "" {
			bookDetails.Genre = updateBook.UploadedCoverImage
		}

		if err := db.Save(&bookDetails).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error updating book: ", err)
		}
		res, _ := json.Marshal(bookDetails)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to edit this book!")
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)
	bookId := params["id"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		log.Println("Error while parsing!: ", err)
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	bookUserId, _ := models.GetBookById(ID)

	if bookUserId.UserId == id {
		book := models.DeleteBook(ID)
		res, _ := json.Marshal(book)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to delete this book!")
	}
}

func DeleteCoverFromS3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	id, _ := utils.JwtUserIdUsername(w, r)
	intId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		log.Println("Failed to parse: ", err)
	}
	book, db := models.GetBookById(intId)

	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete the profile image from S3 (if it exists)
	if book.UploadedCoverImage != "" {
		err := utils.DeleteFromS3(book.CoverImageUrl)
		if err != nil {
			log.Println("Error deleting from S3:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		book.CoverImageUrl = ""      // Clear the aws s3 URL in the database
		book.UploadedCoverImage = "" // Clear the uploaded URL in the database
	}

	// Update the user profile in the database
	if err := db.Save(&book).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error updating user: ", err)
		return
	}

	w.Write([]byte("Cover image deleted successfully"))
}

func DeleteBookFromS3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	id, _ := utils.JwtUserIdUsername(w, r)
	intId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		log.Println("Failed to parse: ", err)
	}
	book, db := models.GetBookById(intId)

	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete the profile image from S3 (if it exists)
	if book.UploadedBook != "" {
		err := utils.DeleteFromS3(book.FileURL)
		if err != nil {
			log.Println("Error deleting from S3:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		book.FileURL = "" // Clear the AvatarURL in the database
		book.UploadedBook = ""
	}

	// Update the user profile in the database
	if err := db.Save(&book).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("Error updating user: ", err)
		return
	}

	w.Write([]byte("Book file deleted successfully"))
}
