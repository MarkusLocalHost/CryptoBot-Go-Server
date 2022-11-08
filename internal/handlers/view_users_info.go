package handlers

import (
	"cryptocurrency/internal/middleware"
	"cryptocurrency/internal/utils/apperrors"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type viewUsersInfoReq struct {
	UserID string `json:"userID"`
}

func (h *Handler) ViewUsersInfo(c *gin.Context) {
	var req viewUsersInfoReq

	// get context data
	middlewareData, _ := c.Get(middleware.AuthorizationPayloadKey)
	data := middlewareData.(string)
	data = strings.ReplaceAll(data, "'", "\"")
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		log.Println(err)
	}

	// log request
	ctx := c.Request.Context()
	err = h.LogService.LogRequestFromSite(ctx, "/manager/info/users_info", req.UserID, req)
	if err != nil {
		log.Fatal(err)
	}

	// save observer to db
	usersInfo, err := h.ManagerService.ViewUsersInfo(ctx)
	if err != nil {
		log.Printf("Failed to create promo code: %v\n", err.Error())
		c.JSON(
			apperrors.Status(err),
			gin.H{"error": err})
		return
	}

	// return result with status
	c.JSON(http.StatusCreated, gin.H{
		"result": usersInfo,
	})
}
