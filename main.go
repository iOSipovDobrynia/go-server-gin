package main

import (
	"context"
	handler "go-server-gin/handler"
	storage "go-server-gin/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	s := storage.NewStorage(pool)
	h := handler.NewHandler(s)

	router := gin.Default()
	router.POST(
		"/addDriver", h.AddDriverHandler,
	)
	router.GET(
		"/driver/:id", h.GetDriverHandler,
	)
	router.GET(
		"/drivers", h.GetDriverListHandler,
	)

	router.Run()
}
