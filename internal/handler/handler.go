package handler

import (
	"context"
	"errors"
	"go-server-gin/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	service Service
}

type Service interface {
	AddDriver(ctx context.Context, newDriver domain.Driver) (domain.AddDriverResponse, error)
	GetDriverById(ctx context.Context, driverID domain.DriverIdRequest) (*domain.Driver, error)
	GetFullDriverById(ctx context.Context, driverID domain.DriverIdRequest) (*domain.FullDriver, error)
	GetDriverList(ctx context.Context) ([]domain.Driver, error)
}

func New(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddDriverHandler(ctx *gin.Context) {
	var newDriver domain.Driver

	if err := ctx.ShouldBindJSON(&newDriver); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	addDriverResponse, err := h.service.AddDriver(ctx.Request.Context(), newDriver)
	if err != nil {
		log.Printf("Error: %v", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "vehicle_id does not exist"})
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": addDriverResponse.ID})
}

func (h *Handler) GetDriverHandler(ctx *gin.Context) {
	var driverID domain.DriverIdRequest

	if err := ctx.ShouldBindUri(&driverID); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	driver, err := h.service.GetDriverById(ctx.Request.Context(), driverID)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, driver)
}

func (h *Handler) GetFullDriverHandler(ctx *gin.Context) {
	var driverID domain.DriverIdRequest

	if err := ctx.ShouldBindUri(&driverID); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	driver, err := h.service.GetFullDriverById(ctx.Request.Context(), driverID)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, driver)
}

func (h *Handler) GetDriverListHandler(ctx *gin.Context) {
	driverList, err := h.service.GetDriverList(ctx.Request.Context())
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, driverList)
}
