package service

import (
	"database/sql"
	"fmt"
	"time"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Genre  string
}

type BookService struct {
	db *sql.DB
}

func NewBookService(db *sql.DB) *BookService {
	return &BookService{db: db}
}

func (s *BookService) CreateBook(book *Book) error {
	query := "INSERT INTO books (tittle, author, genre) VALUES (?,?,?)"
	result, err := s.db.Exec(query, book.Title, book.Author, book.Genre)

	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	book.ID = int(lastInsertID)
	return nil
}

func (s *BookService) GetBooks() ([]Book, error) {
	query := "SELECT id, title, author, genre FROM books"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *BookService) GetBookById(id int) (*Book, error) {
	query := "SELECT id, title, author, genre FROM books WHERE id = ?"
	rows := s.db.QueryRow(query, id)

	var book Book
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *BookService) SearchBooksByName(name string) ([]Book, error) {
	query := "SELECT id, title, author, genre FROM books WHERE title LIKE ?"
	rows, err := s.db.Query(query, "%"+name+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (s *BookService) UpdateBook (book *Book) error {
	query := "UPDATE books SET title = ?, author = ?, genre = ? WHERE id = ?"
	_, err := s.db.Exec(query, book.Title, book.Author, book.Genre, book.ID)
	return err
}

func (s *BookService) DeleteBook (id int) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err	
}

func (s *BookService) SimulateReading(bookID int, duration time.Duration, results chan<- string) {
	book, err := s.GetBookById(bookID)
	if err != nil || book == nil {
		results <- fmt.Sprintf("Book %d not found", bookID)
	}

	time.Sleep(duration)
	results <- fmt.Sprintf("Finished reading %s", book.Title)
}

func (s *BookService) SimulateMultipleReadings(booksIDs []int, duration time.Duration) []string {
	results := make(chan string, len(booksIDs))

	for _, id := range booksIDs {
		go func(bookID int) {
			s.SimulateReading(bookID, duration, results)
		}(id)
	}

	var responses []string
	for range booksIDs {
		responses = append(responses, <-results)
	}
	close(results)
	return responses
}