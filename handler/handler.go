package handler

import (
	"errors"
	storage "go-server-gin/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	storage *storage.Storage
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) AddDriverHandler(ctx *gin.Context) {
	var newDriver storage.Driver

	if err := ctx.ShouldBindJSON(&newDriver); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	addDriverResponse, err := h.storage.AddDriver(ctx.Request.Context(), newDriver)
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
	var driverID storage.DriverIdRequest

	if err := ctx.ShouldBindUri(&driverID); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	driver, err := h.storage.GetDriverById(ctx.Request.Context(), driverID)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, driver)
}

func (h *Handler) GetFullDriverHandler(ctx *gin.Context) {
	var driverID storage.DriverIdRequest

	if err := ctx.ShouldBindUri(&driverID); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	driver, err := h.storage.GetFullDriverById(ctx.Request.Context(), driverID)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, driver)
}

func (h *Handler) GetDriverListHandler(ctx *gin.Context) {
	driverList, err := h.storage.GetDriverList(ctx.Request.Context())
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, driverList)
}
