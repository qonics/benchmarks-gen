package controller

import (
	"github.com/gin-gonic/gin"
)

/*
Receive deleteCache request
*/
func Index(c *gin.Context) {
	//helper.SecurePath(c)
	c.JSON(200, gin.H{"status": 200,
		"message": "Weclome to GIN kickstart project from Qonics inc",
	})
}

func ServiceStatusCheck(c *gin.Context) {
	c.JSON(400, gin.H{"status": 200, "message": "This API service is running"})
}
