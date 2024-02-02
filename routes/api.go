package routes

import (
	"benchmarks-gin/controller"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1/")

	v1.GET("service-status", controller.ServiceStatusCheck)
	v1.POST("dataset/generate", controller.GenerateDataset)

	return r
}
