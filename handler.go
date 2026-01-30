package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage *Storage
}

func NewHandler(storage *Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) AddDriverHandler(ctx *gin.Context) {
	var newDriver Driver

	if err := ctx.ShouldBindJSON(&newDriver); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	addDriverResponse, err := h.storage.AddDriver(ctx.Request.Context(), newDriver)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": addDriverResponse.ID})
}

func (h *Handler) GetDriverHandler(ctx *gin.Context) {
	var driverID DriverIdRequest

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

func (h *Handler) GetDriverListHandler(ctx *gin.Context) {
	driverList, err := h.storage.GetDriverList(ctx.Request.Context())
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, driverList)
}
