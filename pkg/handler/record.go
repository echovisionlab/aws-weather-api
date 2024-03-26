package handler

import (
	"github.com/echovisionlab/aws-weather-api/pkg/query"
	"github.com/echovisionlab/aws-weather-api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func GetRecord(service *service.Service, validate *validator.Validate) func(c *gin.Context) {
	return func(c *gin.Context) {
		var q query.Record
		if err := c.ShouldBindQuery(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if err := validate.Struct(q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		records := service.FindRecordBy(&q)
		c.JSON(http.StatusOK, gin.H{"data": records})
	}
}
