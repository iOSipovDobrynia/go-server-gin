package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Driver struct {
	ID        int64  `json:"id"`
	Name      string `json:"name" binding:"required"`
	VehicleID int64  `json:"vehicle_id" binding:"required"`
	Score     int    `json:"score" binding:"required"`
}

type FullDriver struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Vehicle Vehicle `json:"vehicle"`
	Score   int     `json:"score"`
}

type Vehicle struct {
	ID     int64  `json:"id"`
	Type   string `json:"type"`
	Vendor string `json:"vendor"`
	Model  string `json:"model"`
}

type DriverIdRequest struct {
	ID int64 `json:"id" uri:"id" binding:"required"`
}

type AddDriverResponse struct {
	ID int64 `json:"id"`
}

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) AddDriver(ctx context.Context, newDriver Driver) (AddDriverResponse, error) {
	var addDriverResponse AddDriverResponse
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO drivers (name, vehicle_id, score) VALUES ($1, $2, $3) RETURNING id;",
		newDriver.Name,
		newDriver.VehicleID,
		newDriver.Score,
	).Scan(&addDriverResponse.ID)
	if err != nil {
		return AddDriverResponse{}, err
	}

	return addDriverResponse, nil
}

func (s *Storage) GetDriverById(ctx context.Context, driverID DriverIdRequest) (*Driver, error) {
	var driver Driver
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, name, vehicle_id, score FROM drivers WHERE id=$1;",
		driverID.ID,
	).Scan(&driver.ID, &driver.Name, &driver.VehicleID, &driver.Score)
	if err != nil {
		return nil, err
	}

	return &driver, nil
}

func (s *Storage) GetFullDriverById(ctx context.Context, driverID DriverIdRequest) (*FullDriver, error) {
	var driver FullDriver
	err := s.pool.QueryRow(
		ctx,
		`
			SELECT d.id, d.name, v.id, v.type, v.vendor, v.model, d.score 
			FROM drivers d 
			INNER JOIN vehicles v 
				ON d.vehicle_id = v.id 
			WHERE d.id = $1;
			`,
		driverID.ID,
	).Scan(
		&driver.ID, &driver.Name, &driver.Vehicle.ID, &driver.Vehicle.Type, &driver.Vehicle.Vendor,
		&driver.Vehicle.Model,
		&driver.Score,
	)
	if err != nil {
		return nil, err
	}

	return &driver, nil
}

func (s *Storage) GetDriverList(ctx context.Context) ([]Driver, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, vehicle_id, score FROM drivers;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var driverList []Driver
	for rows.Next() {
		var driver Driver
		err := rows.Scan(&driver.ID, &driver.Name, &driver.VehicleID, &driver.Score)
		if err != nil {
			return nil, err
		}
		driverList = append(driverList, driver)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return driverList, nil
}
