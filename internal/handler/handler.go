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

// AddDriverHandler godoc
// @Summary      Добавить нового водителя
// @Description  Создает запись о водителе в базе данных
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        input   body      domain.Driver  true  "Данные водителя"
// @Success      201     {object}  domain.AddDriverResponse
// @Failure      400     {string}  string "Bad Request"
// @Router       /addDriver [post]
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

// GetDriverHandler godoc
// @Summary      Получить информацию о водителе
// @Description  Получает информацию о водителе по его уникальному идентификатору
// @Tags         drivers
// @Produce      json
// @Param        id   path int true  "ID водителя"
// @Success      200     {object}  domain.Driver
// @Failure      400  {string}  string "Неверный ID"
// @Router       /driver/{id} [get]
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

// GetFullDriverHandler godoc
// @Summary      Получить информацию о водителе и транспортном средстве
// @Description  Получить информацию о водителе и транспортном средстве по его уникальному идентификатору водителя
// @Tags         drivers
// @Produce      json
// @Param        id   path int true  "ID водителя"
// @Success      200     {object}  domain.FullDriver
// @Failure      400  {string}  string "Неверный ID"
// @Router       /driver/{id}/full [get]
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

// GetDriverListHandler godoc
// @Summary      Получить информацию о водителях
// @Description  Получить информацию о водителях
// @Tags         drivers
// @Produce      json
// @Success      200     {array}  domain.Driver
// @Router       /drivers [get]
func (h *Handler) GetDriverListHandler(ctx *gin.Context) {
	driverList, err := h.service.GetDriverList(ctx.Request.Context())
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, driverList)
}
