package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Driver struct {
	ID      int64  `json:"id"`
	Name    string `json:"name" binding:"required"`
	Vehicle string `json:"vehicle" binding:"required"`
	Score   int    `json:"score" binding:"required"`
}

type DriverIdRequest struct {
	ID int64 `json:"id" uri:"id" binding:"required"`
}

type AddDriverResponse struct {
	ID int64 `json:"id"`
}

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

	storage := NewStorage(pool)
	handler := NewHandler(storage)

	router := gin.Default()
	router.POST(
		"/addDriver", handler.AddDriverHandler,
	)
	router.GET(
		"/driver/:id", handler.GetDriverHandler,
	)
	router.GET(
		"/drivers", handler.GetDriverListHandler,
	)

	router.Run()
}
