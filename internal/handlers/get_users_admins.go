package handlers

import (
	"cryptocurrency/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type getUsersAdminsReq struct{}

// GetUsersAdmins Get currency which enable to track
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
func (h *Handler) GetUsersAdmins(c *gin.Context) {
	var req getUsersAdminsReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	err := json.Unmarshal([]byte(middlewareData.(string)), &req)
	if err != nil {
		log.Printf("Failed to parse request body: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// save observer to db
	ctx := c.Request.Context()
	data, err := h.InfoService.GetAllUsersAdmins(ctx)
	if err != nil {
		log.Printf("Failed to fetch admins: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}
