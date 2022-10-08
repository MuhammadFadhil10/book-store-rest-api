package controller

import (
	"encoding/json"
	"fmt"
	"go-rest/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAllBook(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type","application/json");
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	books := models.Book.FetchAllBooks(models.Book{})
	resp := map[string]interface{} {
		"status": true,
		"totalresult": len(books),
		"data": books,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println(err)
	}
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	userId := mux.Vars(r)["userId"]
	var b models.Book
	decodeErr := json.NewDecoder(r.Body).Decode(&b)
	if decodeErr != nil {
		fmt.Println(decodeErr)
	}
	b.Id = uuid.New().String()

	insertErr := models.Book.AddBook(models.Book{},b.Id,b.Title,b.Page,b.Year,b.Genre,b.Synopsis,b.Price,b.Quantity,b.Image,userId)
	resp := make(map[string]interface{})
	if insertErr != nil {
		w.WriteHeader(500)
		return;
	}
	resp["status"] = true
	resp["message"] = "Cool, Your Book Are Published!"
	json.NewEncoder(w).Encode(resp)

}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	bookId := mux.Vars(r)["bookId"]

	var updatedBook models.Book

	decodeErr := json.NewDecoder(r.Body).Decode(&updatedBook)

	if decodeErr != nil {
		log.Fatal(decodeErr)
	}

	insertErr := models.Book.UpdateBook(
		models.Book{},updatedBook.Title,updatedBook.Page,updatedBook.Year,updatedBook.Genre,updatedBook.Synopsis,
		updatedBook.Price,updatedBook.Quantity,updatedBook.Image,bookId,
	)
	if insertErr != nil {
		w.WriteHeader(500)
		log.Fatal(insertErr)
		return
	}

	resp := map[string]interface{} {
		"status": true,
		"message": "Your book has Updated!",
	}
	json.NewEncoder(w).Encode(resp)

}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	bookId := mux.Vars(r)["bookId"]

	deleteErr := models.Book.DeleteBook(models.Book{},bookId)

	if deleteErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp := map[string]interface{} {
		"status": true,
		"message": "Book has been Deleted!",
	}
	json.NewEncoder(w).Encode(resp)

}