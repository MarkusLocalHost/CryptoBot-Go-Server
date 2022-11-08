package handlers

import (
	"cryptocurrency/internal/middleware"
	"cryptocurrency/internal/utils/apperrors"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type getUsersLanguagesReq struct{}

// GetUsersLanguages Get currency which enable to track
// @summary View supported vs currencies
// @description View currency which enable to track
// @Security Bearer
// @Tags Info
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @router /info/supported_vs_currencies/view [get]
func (h *Handler) GetUsersLanguages(c *gin.Context) {
	var req getUsersLanguagesReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	err := json.Unmarshal([]byte(middlewareData.(string)), &req)
	if err != nil {
		log.Println(err)
	}

	// save observer to db
	ctx := c.Request.Context()
	data, err := h.InfoService.GetAllUsersLanguages(ctx)
	if err != nil {
		log.Printf("Failed to fetch data: %v\n", err.Error())
		c.JSON(
			apperrors.Status(err),
			gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}
