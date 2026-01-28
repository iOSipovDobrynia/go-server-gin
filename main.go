package main

import (
	"context"
	"log"
	"net/http"
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
	ID int64
}

var pool *pgxpool.Pool

func AddDriverHandler(c *gin.Context) {
	var newDriver Driver

	if err := c.ShouldBindJSON(&newDriver); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var addDriverResponse AddDriverResponse
	err := pool.QueryRow(
		c.Request.Context(),
		"INSERT INTO drivers (name, vehicle, score) VALUES ($1, $2, $3) RETURNING id;",
		newDriver.Name,
		newDriver.Vehicle,
		newDriver.Score,
	).Scan(&addDriverResponse.ID)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": addDriverResponse.ID})
}

func GetDriverHandler(c *gin.Context) {
	var driverID DriverIdRequest

	if err := c.ShouldBindUri(&driverID); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var driver Driver
	err := pool.QueryRow(
		c.Request.Context(),
		"SELECT id, name, vehicle, score FROM drivers WHERE id=$1",
		driverID.ID,
	).Scan(&driver.ID, &driver.Name, &driver.Vehicle, &driver.Score)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, driver)
}

func GetDriverListHandler(c *gin.Context) {
	rows, err := pool.Query(c.Request.Context(), "SELECT id, name, vehicle, score FROM drivers")
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	defer rows.Close()

	var driverList []Driver
	for rows.Next() {
		var driver Driver
		err := rows.Scan(&driver.ID, &driver.Name, &driver.Vehicle, &driver.Score)
		if err != nil {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
		driverList = append(driverList, driver)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, driverList)
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

	pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	router := gin.Default()
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
