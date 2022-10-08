package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn
var err error

func ConnectDb() {
	
	Conn,err = pgx.Connect(context.Background(),"postgres://postgres:1010@localhost:5432/db_bookstore")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error connecting to database")
		return
	}

	fmt.Println("Database Connected!")
}