package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
		"INSERT INTO drivers (name, vehicle, score) VALUES ($1, $2, $3) RETURNING id;",
		newDriver.Name,
		newDriver.Vehicle,
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
		"SELECT id, name, vehicle, score FROM drivers WHERE id=$1;",
		driverID.ID,
	).Scan(&driver.ID, &driver.Name, &driver.Vehicle, &driver.Score)
	if err != nil {
		return nil, err
	}

	return &driver, nil
}

func (s *Storage) GetDriverList(ctx context.Context) ([]Driver, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, vehicle, score FROM drivers;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var driverList []Driver
	for rows.Next() {
		var driver Driver
		err := rows.Scan(&driver.ID, &driver.Name, &driver.Vehicle, &driver.Score)
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
