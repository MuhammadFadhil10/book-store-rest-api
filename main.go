package main

import (
	"fmt"
	"go-rest/controller"
	"go-rest/db"
	"go-rest/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter();
	db.ConnectDb()	
	
	
	// subRouter := r.PathPrefix("/").Subrouter()
	r.HandleFunc("/books", middleware.JwtMiddleware(controller.GetAllBook)).Methods("GET")
	r.HandleFunc("/add-book/{userId}", middleware.JwtMiddleware(controller.AddBook)).Methods("POST")
	r.HandleFunc("/update-book/{bookId}", middleware.JwtMiddleware(controller.UpdateBook)).Methods("PUT")
	r.HandleFunc("/delete-book/{bookId}", middleware.JwtMiddleware(controller.DeleteBook)).Methods("DELETE")
	r.HandleFunc("/register", controller.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", controller.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", middleware.JwtMiddleware(controller.Logout)).Methods("GET")

	
	fmt.Println("server run at port 8000")
	listenErr := http.ListenAndServe(":8000",r)
	if listenErr != nil {
		log.Fatal(listenErr)
	}


}

