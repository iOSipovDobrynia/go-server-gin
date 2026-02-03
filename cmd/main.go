package main

import (
	"context"
	"go-server-gin/internal/handler"
	"go-server-gin/internal/service"
	"go-server-gin/internal/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	_ "go-server-gin/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title driver api
// @version         1.0
// @description     Это сервер для управления водителями.
// @host            localhost:8080
// @BasePath        /
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

	s := storage.New(pool)
	svc := service.New(s)
	h := handler.New(svc)

	router := gin.Default()
	router.POST(
		"/addDriver", h.AddDriverHandler,
	)
	router.GET(
		"/driver/:id", h.GetDriverHandler,
	)
	router.GET(
		"/driver/:id/full", h.GetFullDriverHandler,
	)
	router.GET(
		"/drivers", h.GetDriverListHandler,
	)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
