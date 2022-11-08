package bind_data

import (
	"cryptocurrency/internal/utils/apperrors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func BindData(c *gin.Context, req interface{}) bool {
	// send error if Content-Type != application/json
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())

		err := apperrors.NewUnsupportedMediaType(msg)

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return false
	}

	if err := c.ShouldBind(req); err != nil {
		log.Printf("Error binding data %+v\n", err)

		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}

			err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

			c.JSON(err.Status(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}

		fallBack := apperrors.NewInternal()

		c.JSON(fallBack.Status(), gin.H{
			"error": fallBack,
		})
		return false
	}

	return true
}
