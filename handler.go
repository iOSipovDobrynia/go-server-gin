package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{storage: storage}
}

func (handler *Handler) AddDriverHandler(ctx *gin.Context) {
	var newDriver Driver

	if err := ctx.ShouldBindJSON(&newDriver); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	addDriverResponse, err := handler.storage.AddDriver(ctx, newDriver)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": addDriverResponse.ID})
}

func (handler *Handler) GetDriverHandler(ctx *gin.Context) {
	var driverID DriverIdRequest

	if err := ctx.ShouldBindUri(&driverID); err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	driver, err := handler.storage.GetDriverById(ctx, driverID)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, driver)
}

func (handler *Handler) GetDriverListHandler(ctx *gin.Context) {
	driverList, err := handler.storage.GetDriverList(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, driverList)
}
