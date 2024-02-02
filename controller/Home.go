package controller

import (
	"benchmarks-gin/config"
	"benchmarks-gin/helper"
	"benchmarks-gin/model"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateDataset(c *gin.Context) {
	startTime := time.Now()
	payload := model.Payload{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(417, gin.H{"status": 417,
			"message": "Please provide all the required data", "details": err.Error()})
		return
	}
	if len(payload.Database) == 0 {
		c.JSON(417, gin.H{"status": 417,
			"message": "Please provide the database for parent data"})
		return
	}
	payload.Result = payload.Records
	config.ConnectDb(payload.Database)
	// defer config.SESSION.Close()
	// defer config.DB.Close()
	// var wg sync.WaitGroup
	// const maxRoutines = 20
	// waitingChannel := make(chan int, maxRoutines)
	i := 0
	for i < int(payload.RecordCount) {
		result := make(map[string]json.RawMessage)
		increment, err := helper.GenerateSqlQuery(payload, result, nil)
		if err != nil {
			c.JSON(417, gin.H{"status": 417,
				"message": "Dataset generation failed", "details": err.Error()})
			return
		}
		i = i + increment
		var failed []map[string]json.RawMessage
		for _, children := range payload.Children {
			children.Result = make(map[string]json.RawMessage)
			ii := 0
			for ii < int(children.RecordCount) {
				increment, err := helper.GenerateSqlQuery(children, children.Result, result)
				if err != nil {
					failed = append(failed, map[string]json.RawMessage{"item": children.Records["id"], "table": json.RawMessage(children.Table), "err": json.RawMessage(err.Error())})
					increment = 1
				}
				ii = ii + increment
			}
		}
		fmt.Println("*********************** Loop end here ************************", i)
	}
	// wg.Wait()
	endTime := time.Now()
	diff := endTime.Sub(startTime)

	c.JSON(200, gin.H{"status": 200,
		"message": fmt.Sprintf("Dataset generated, started at %v, ends at %v. Minutes taken %v", startTime.Format(time.TimeOnly), endTime.Format(time.TimeOnly), diff.Minutes()),
	})
}

func ServiceStatusCheck(c *gin.Context) {
	c.JSON(400, gin.H{"status": 200, "message": "This API service is running"})
}
