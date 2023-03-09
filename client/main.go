package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"net/http"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	baseURL  = "http://localhost:8080"
	speedCof = 100
)

func main() {
	bu := os.Getenv("BASE_URL")
	if bu != "" {
		baseURL = bu
	}

	sc := os.Getenv("SPEED_COF")
	if sc != "" {
		if parsed, err := strconv.ParseInt(sc, 0, 32); err != nil {
			speedCof = int(parsed)
		}
	}

	rand.Seed(time.Now().UnixNano())
	for {
		d := rand.Intn(3*speedCof) + 1*speedCof
		time.Sleep(time.Millisecond * time.Duration(d))
		n := rand.Intn(100)
		switch {
		case n < 10:
			go getAll(n)
		case n >= 10 && n < 30:
			go getSingleBook(n)
		case n >= 30 && n < 50:
			go creteNew(n)
		case n >= 50 && n < 75:
			go updateBook(n)
		case n >= 75:
			go deleteBook(n)
		}
	}
}

// Get all books
func getAll(n int) {
	req, _ := http.NewRequest("GET", baseURL+"/books", nil)
	req.Header.Set("TraceID", fmt.Sprint(n))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var books []Book
	if err := json.Unmarshal(body, &books); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All books:")
	for _, book := range books {
		fmt.Printf("- %s by %s\n", book.Title, book.Author)
	}
}

func getSingleBook(n int) {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(30)
	bookID := fmt.Sprint(i)

	req, _ := http.NewRequest("GET", baseURL+"/book/"+bookID, nil)
	req.Header.Set("TraceID", fmt.Sprint(n))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var book Book
	if err := json.Unmarshal(body, &book); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nBook with ID %s:\n", bookID)
	fmt.Printf("- %s by %s\n", book.Title, book.Author)
}

func creteNew(n int) {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(30)
	bookID := fmt.Sprint(i)
	// Create new book
	newBook := Book{ID: bookID, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"}
	jsonBook, err := json.Marshal(newBook)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, _ := http.NewRequest("POST", baseURL+"/book", bytes.NewBuffer(jsonBook))
	req.Header.Set("TraceID", fmt.Sprint(n))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var createdBook Book
	if err := json.Unmarshal(body, &createdBook); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nCreated book:\n")
	fmt.Printf("- %s by %s\n", createdBook.Title, createdBook.Author)
}

func updateBook(n int) {
	// Update book
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(30)
	bookID := fmt.Sprint(i)
	updateBook := Book{ID: bookID, Title: "The Great Gatsby (Updated)", Author: "F. Scott Fitzgerald"}
	jsonUpdateBook, err := json.Marshal(updateBook)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", baseURL+"/book/"+updateBook.ID, bytes.NewBuffer(jsonUpdateBook))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TraceID", fmt.Sprint(n))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var updatedBook Book
	if err := json.Unmarshal(body, &updatedBook); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\nUpdated book:\n")
	fmt.Printf("- %s by %s\n", updatedBook.Title, updatedBook.Author)
}

func deleteBook(n int) {
	// Delete book
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(30)
	bookID := fmt.Sprint(i)
	req, err := http.NewRequest("DELETE", baseURL+"/book/"+bookID, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("TraceID", fmt.Sprint(n))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("\nDeleted book with ID %s\n", bookID)
}
