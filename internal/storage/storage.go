package storage

import (
	"context"
	"go-server-gin/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) AddDriver(ctx context.Context, newDriver domain.Driver) (domain.AddDriverResponse, error) {
	var addDriverResponse domain.AddDriverResponse
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO drivers (name, vehicle_id, score) VALUES ($1, $2, $3) RETURNING id;",
		newDriver.Name,
		newDriver.VehicleID,
		newDriver.Score,
	).Scan(&addDriverResponse.ID)
	if err != nil {
		return domain.AddDriverResponse{}, err
	}

	return addDriverResponse, nil
}

func (s *Storage) GetDriverById(ctx context.Context, driverID domain.DriverIdRequest) (*domain.Driver, error) {
	var driver domain.Driver
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

func (s *Storage) GetFullDriverById(ctx context.Context, driverID domain.DriverIdRequest) (*domain.FullDriver, error) {
	var driver domain.FullDriver
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

func (s *Storage) GetDriverList(ctx context.Context) ([]domain.Driver, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, vehicle_id, score FROM drivers;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var driverList []domain.Driver
	for rows.Next() {
		var driver domain.Driver
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
