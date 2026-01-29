package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) Storage {
	return Storage{pool: pool}
}

func (s *Storage) AddDriver(ctx *gin.Context, newDriver Driver) (AddDriverResponse, error) {
	var addDriverResponse AddDriverResponse
	err := s.pool.QueryRow(
		ctx.Request.Context(),
		"INSERT INTO drivers (name, vehicle, score) VALUES ($1, $2, $3) RETURNING id;",
		newDriver.Name,
		newDriver.Vehicle,
		newDriver.Score,
	).Scan(&addDriverResponse.ID)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return AddDriverResponse{}, err
	}

	return addDriverResponse, nil
}

func (s *Storage) GetDriverById(ctx *gin.Context, driverID DriverIdRequest) (Driver, error) {
	var driver Driver
	err := s.pool.QueryRow(
		ctx.Request.Context(),
		"SELECT id, name, vehicle, score FROM drivers WHERE id=$1;",
		driverID.ID,
	).Scan(&driver.ID, &driver.Name, &driver.Vehicle, &driver.Score)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return Driver{}, err
	}

	return driver, nil
}

func (s *Storage) GetDriverList(ctx *gin.Context) ([]Driver, error) {
	rows, err := s.pool.Query(ctx.Request.Context(), "SELECT id, name, vehicle, score FROM drivers;")
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return []Driver{}, err
	}
	defer rows.Close()

	var driverList []Driver
	for rows.Next() {
		var driver Driver
		err := rows.Scan(&driver.ID, &driver.Name, &driver.Vehicle, &driver.Score)
		if err != nil {
			log.Printf("Error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return []Driver{}, err
		}
		driverList = append(driverList, driver)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return []Driver{}, err
	}

	return driverList, nil
}
