package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matthinc/gomment/logic"
)

func handleValidationError(c *gin.Context, err logic.ValidationError) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error_type":            "validation_error",
		"validation_field_name": err.FieldName,
		"validation_type":       err.ValidationType,
		"validation_info":       err.Info,
	})
}
