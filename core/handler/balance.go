package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/whyslove/avito-test/core/models"
)

func (h *Handler) GetBalance(c *gin.Context) {
	id := c.Param("id") //always have in other way we will not fall into this route
	logrus.Debugf("Receive get balance request with id=%s", id)
	numeric_id, err := strconv.Atoi(id)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	balance, err := h.services.GetBalance(numeric_id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if currency := c.Query("currency"); currency != "" {
		balance, err = h.currencyConverter.ConvertFromRub(currency, balance)
		if err != nil {
			logrus.Error(err.Error())
			NewErrorResponse(c, http.StatusInternalServerError, "Error in currency converting")
			return
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"balance": balance,
	})

}

func (h *Handler) OperationBalance(c *gin.Context) {
	var operation models.UserOperation
	var err error

	if err := c.BindJSON(&operation); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if currency := c.Query("currency"); currency != "" {
		operation.AmountOfMoney, err =
			h.currencyConverter.ConvertToRub(currency, operation.AmountOfMoney)
		if err != nil {
			logrus.Error(err.Error())
			NewErrorResponse(c, http.StatusInternalServerError, "Error in currency converting")
			return
		}
	}
	if operation.AmountOfMoney < 0 {
		err = h.services.WriteOffMoney(operation)
	} else {
		err = h.services.AddMoney(operation)
	}
	if err != nil {
		if errors.Is(err, models.NoRecordInDb) {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) TransferBalance(c *gin.Context) {
	var userTransferOperation models.UserTransferOperation
	var err error

	if err := c.BindJSON(&userTransferOperation); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if currency := c.Query("currency"); currency != "" {
		userTransferOperation.AmountOfMoney, err =
			h.currencyConverter.ConvertToRub(currency, userTransferOperation.AmountOfMoney)
		if err != nil {
			logrus.Error(err.Error())
			NewErrorResponse(c, http.StatusInternalServerError, "Error in currency converting")
			return
		}
	}
	err = h.services.TransferOperation(userTransferOperation)
	if err != nil {
		if errors.Is(err, models.NoRecordInDb) {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
