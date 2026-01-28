package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Driver struct {
	ID      int64  `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Vehicle string `json:"vehicle" binding:"required"`
	Score   int    `json:"score" binding:"required"`
}

type DriverIdRequest struct {
	ID int64 `json:"id" uri:"id" binding:"required"`
}

var pool *pgxpool.Pool

func PingHandler(c *gin.Context) {
	c.JSON(
		http.StatusOK, gin.H{
			"message": "pong",
		},
	)
}

func AddDriverHandler(c *gin.Context) {
	var newDriver Driver

	if errBind := c.ShouldBindJSON(&newDriver); errBind != nil {
		c.JSON(http.StatusBadRequest, errBind.Error())
		return
	}

	_, errDb := pool.Exec(
		c.Request.Context(),
		"INSERT INTO drivers (id, name, vehicle, score) VALUES ($1, $2, $3, $4)",
		newDriver.ID,
		newDriver.Name,
		newDriver.Vehicle,
		newDriver.Score,
	)
	if errDb != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not add driver"})
		return
	}

	c.JSON(http.StatusCreated, newDriver)
}

func GetDriverHandler(c *gin.Context) {
	var driverID DriverIdRequest

	if errBindUri := c.ShouldBindUri(&driverID); errBindUri != nil {
		c.JSON(http.StatusBadRequest, errBindUri.Error())
		return
	}

	var driver Driver
	errDb := pool.QueryRow(
		c.Request.Context(),
		"SELECT id, name, vehicle, score FROM drivers WHERE id=$1",
		driverID.ID,
	).Scan(&driver.ID, &driver.Name, &driver.Vehicle, &driver.Score)

	if errors.Is(errDb, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "driver not found"})
		return
	} else if errDb != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}

	c.JSON(http.StatusOK, driver)
}

func GetDriverListHandler(c *gin.Context) {
	var driverList []Driver

	rows, errDb := pool.Query(c.Request.Context(), "SELECT id, name, vehicle, score FROM drivers")

	if errDb != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var driver Driver
		errScan := rows.Scan(&driver.ID, &driver.Name, &driver.Vehicle, &driver.Score)
		if errScan != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "scan error"})
			return
		}
		driverList = append(driverList, driver)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "rows error"})
		return
	}

	c.JSON(http.StatusOK, driverList)
}

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Выполняем запрос: считаем количество пользователей
	row := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM drivers")
	var count int
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Пользователей в системе:", count)

	router := gin.Default()
	router.GET(
		"/ping", PingHandler,
	)
	router.POST(
		"/addDriver", AddDriverHandler,
	)
	router.GET(
		"/driver/:id", GetDriverHandler,
	)
	router.GET(
		"/drivers", GetDriverListHandler,
	)

	router.Run()
}
