package service

import (
	"context"
	"go-server-gin/internal/domain"
)

type Service struct {
	storage DriverStorage
}

type DriverStorage interface {
	AddDriver(ctx context.Context, newDriver domain.Driver) (domain.AddDriverResponse, error)
	GetDriverById(ctx context.Context, driverID domain.DriverIdRequest) (*domain.Driver, error)
	GetFullDriverById(ctx context.Context, driverID domain.DriverIdRequest) (*domain.FullDriver, error)
	GetDriverList(ctx context.Context) ([]domain.Driver, error)
}

func NewService(storage DriverStorage) *Service {
	return &Service{storage: storage}
}

func (s *Service) AddDriver(ctx context.Context, newDriver domain.Driver) (domain.AddDriverResponse, error) {
	response, err := s.storage.AddDriver(ctx, newDriver)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *Service) GetDriverById(ctx context.Context, driverID domain.DriverIdRequest) (*domain.Driver, error) {
	response, err := s.storage.GetDriverById(ctx, driverID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) GetFullDriverById(ctx context.Context, driverID domain.DriverIdRequest) (
	*domain.FullDriver,
	error,
) {
	response, err := s.storage.GetFullDriverById(ctx, driverID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) GetDriverList(ctx context.Context) ([]domain.Driver, error) {
	response, err := s.storage.GetDriverList(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}
