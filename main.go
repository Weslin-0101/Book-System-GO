package main

import (
	"errors"
	"fmt"
	"gobooks/internal/service"
)

func sumAndCount(x, y int) (int, error) {
	if (x < 0 || y < 0) {
		return 0, errors.New("x and y must be greater than 0")
	}

	return x + y, nil
}

func main() {
	book := service.Book {
		ID: 1,
		Title: "The Hobbit",
		Author: "J.R.R. Tolkien",
		Genre: "Fantasy",
	}

	fmt.Println(book.GetFullBook())
}