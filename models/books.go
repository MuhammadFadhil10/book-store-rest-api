package models

import (
	"context"
	"go-rest/db"
	"log"
)

type Book struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Synopsis string `json:"synopsis"`
	Image string `json:"image"`
	UserId string `json:"userid"`
	Page int `json:"page"`
	Year int `json:"year"`
	Price int `json:"price"`
	Quantity        int `json:"quantity"`
	Genre                              []string `json:"genre"`
	Author Author `json:"author"`
}

type Author struct {
	AuthorId string `json:"author_id"`
	AuthorName string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	AuthorProfilePict string `json:"author_profile_pict"`
	AuthorBio string `json:"author_biography"`
}

type get interface {
	FetchAllBooks() []Book
	FindOneBook(bookId string) Book
}

type post interface {
	AddBook(id string,title string, page int, year int, genre []string, synopsis string, price int, quantity int, image string, userId string) error
}

type put interface {
	UpdateBook(title string, page int, year int, genre []string, synopsis string, price int, quantity int, image string, bookId string) error
}

type delete interface {
	DeleteBook(bookId string) error
}

func (book Book) FetchAllBooks() []Book {
	Conn := db.Conn
	var books []Book
	queryString := `
	SELECT books.id,books.title,books.page,books.year,books.genre,books.synopsis,books.price,
	books.quantity,books.image, users.id as authorid, users.name as authorname, users.email as authoremail,
	users.profile_picture as authorprofilepict, users.biography as authorbiography FROM public.books LEFT JOIN users ON books.user_id = users.id
`	
	data, err := Conn.Query(context.Background(),queryString)
	if err != nil {
		log.Fatal(err)
	}

	for data.Next() {
		
		scanErr := data.Scan(&book.Id,&book.Title,&book.Page, &book.Year, &book.Genre,&book.Synopsis, &book.Price,&book.Quantity, &book.Image,
			&book.Author.AuthorId, &book.Author.AuthorName, &book.Author.AuthorEmail,&book.Author.AuthorProfilePict,
			&book.Author.AuthorBio,
		);
		if scanErr != nil {
			log.Fatal(scanErr)
		}
		books = append(books, book)
	}
	return books
}

func (book Book) AddBook(id string,title string, page int, year int, genre []string, synopsis string, price int, quantity int, image string, userId string) error {

	queryString := `
		INSERT INTO books(id,title,page,year,genre,synopsis,price,quantity,image,user_id)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	Conn := db.Conn

	_, insertErr := Conn.Exec(context.Background(),queryString,id,title,page,year,genre,synopsis,price,quantity,image,userId)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (book Book) UpdateBook(title string, page int, year int, genre []string, synopsis string, price int, quantity int, image string, bookId string) error {

	queryString := `
		UPDATE public.books
		SET title=$1, page=$2, year=$3, genre=array_cat(genre, $4), synopsis=$5, price=$6, quantity=$7, image=$8, updated_at=now()
		WHERE id = $9;
	`
	Conn := db.Conn
		_, updateErr := Conn.Exec(context.Background(),queryString,title,page,year,genre,synopsis,price,quantity,image,bookId)

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (book Book) DeleteBook(bookId string) error {
	Conn := db.Conn

	_, deleteErr := Conn.Exec(context.Background(),"DELETE FROM public.books WHERE id = $1",bookId)

	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

