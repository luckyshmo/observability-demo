package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/uptrace/opentelemetry-go-extra/otelplay"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

type Book struct {
	ID     string `json:"id" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

var books = []Book{
	{ID: "1", Title: "Book 1", Author: "Author 1"},
	{ID: "2", Title: "Book 2", Author: "Author 2"},
	{ID: "3", Title: "Book 3", Author: "Author 3"},
}

func main() {
	router := gin.New()

	// Tracing
	ctx := context.Background()
	shutdown := otelplay.ConfigureOpentelemetry(ctx)
	defer shutdown()

	router.Use(otelgin.Middleware("book-store"))

	// Async logging
	wr := diode.NewWriter(os.Stdout, 1000, 10*time.Millisecond, func(missed int) {
		fmt.Printf("Logger Dropped %d messages\n", missed)
	})
	log := zerolog.New(wr)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Metrics init
	router.GET("/metrics", prometheusHandler())

	router.Use(func(c *gin.Context) {
		now := time.Now()
		r := rand.Intn(50)
		if r == 25 {
			rt := rand.Intn(50) + 50
			time.Sleep(time.Millisecond * time.Duration(rt))
		}

		c.Next()

		ctx := c.Request.Context()
		span := trace.SpanFromContext(ctx)
		id := span.SpanContext().TraceID().String()

		reqByURI.WithLabelValues(c.FullPath(), c.Request.Method, strconv.Itoa(c.Writer.Status())).Inc()
		reqDuration.(prometheus.ExemplarObserver).ObserveWithExemplar(
			time.Since(now).Seconds(), prometheus.Labels{"TraceID": id})
	})

	// Get all books
	router.GET("/books", func(c *gin.Context) {
		log.Info().
			Int("books", len(books)).
			Msg("Books in store left")

		booksCount.Set(float64(len(books)))

		c.JSON(http.StatusOK, books)
	})

	// Get single book
	router.GET("/book/:id", func(c *gin.Context) {
		id := c.Param("id")
		var book Book
		for _, b := range books {
			if b.ID == id {
				book = b
				break
			}
		}

		if book.ID == "" {
			log.Warn().Msg("book not found")
			c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
			return
		}

		log.Info().
			Str("book_id", book.ID).
			Msg("get this ID")

		c.JSON(http.StatusOK, book)
	})

	// Create new book
	router.POST("/book", func(c *gin.Context) {
		log.Error().Err(fmt.Errorf("test post error")).Msg("test lol")

		var book Book
		if err := c.ShouldBindJSON(&book); err != nil {
			log.Error().Err(err).Msg("post book")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Info().
			Str("book_title", book.Title).
			Msg("book title")

		books = append(books, book)
		c.JSON(http.StatusOK, book)
	})

	// Update book
	router.PUT("/book/:id", func(c *gin.Context) {
		log.Error().Err(fmt.Errorf("test kek")).Msg("test err")

		id := c.Param("id")
		var book Book
		for _, b := range books {
			if b.ID == id {
				book = b
				break
			}
		}

		if book.ID == "" {
			log.Error().Int("lol_field", -1).Err(fmt.Errorf("not found")).Msg("put book")
			c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
			return
		}

		if err := c.ShouldBindJSON(&book); err != nil {
			log.Error().Err(err).Msg("put book")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for i, b := range books {
			if b.ID == id {
				books[i] = book
				break
			}
		}

		c.JSON(http.StatusOK, book)
	})
	// Delete book
	router.DELETE("/book/:id", func(c *gin.Context) {
		id := c.Param("id")
		var book Book
		for _, b := range books {
			if b.ID == id {
				book = b
				break
			}
		}

		if book.ID == "" {
			log.Warn().Msg("book not found")
			c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
			return
		}

		for i, b := range books {
			if b.ID == id {
				books = append(books[:i], books[i+1:]...)
				break
			}
		}

		log.Info().
			Str("book_id", book.ID).
			Msg("book deleted")

		c.JSON(http.StatusOK, book)
	})

	// Start server
	router.Run(":80")
}
