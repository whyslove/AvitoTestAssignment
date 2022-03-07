package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whyslove/avito-test/core/models"
)

func (h *Handler) GetOperations(c *gin.Context) {
	var err error
	var strPage string

	id := c.Param("id") //always have otherways we will not fall into this route
	numeric_id, err := strconv.Atoi(id)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if strPage = c.Query("page"); strPage == "" {
		strPage = "1"
	}
	page, err := strconv.Atoi(strPage)

	if err != nil || page < 1 {
		NewErrorResponse(c, http.StatusBadRequest, "Bad page number")
		return
	}
	var operations []models.Operation
	if sort := c.Query("sort"); sort == "" {
		operations, err = h.services.GetOperations(numeric_id, page)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		operations, err = h.services.GetSortedOperations(numeric_id, page, sort)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, operations)

}
