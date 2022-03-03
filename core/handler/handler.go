package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/whyslove/avito-test/core/service"
	"github.com/whyslove/avito-test/pkg/converter"
)

type Handler struct {
	services          *service.Service
	currencyConverter *converter.CurrencyConverter
}

func NewHandler(service *service.Service, currencyConverter *converter.CurrencyConverter) *Handler {
	return &Handler{services: service, currencyConverter: currencyConverter}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("balance/get/:id", h.GetBalance)
		api.POST("balance/operation", h.OperationBalance)
		api.POST("balance/transfer", h.TransferBalance)
	}

	return router
}
