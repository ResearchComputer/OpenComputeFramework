package server

import (
	"fmt"
	"ocfcore/internal/common/structs"
	"ocfcore/internal/server/queue"

	"github.com/gin-gonic/gin"
)

func InferenceRequest(c *gin.Context) {
	var request structs.InferenceStruct
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	topic := "inference:" + request.UniqueModelName
	msg, err := queue.Publish(topic, []byte(fmt.Sprintf("%s", request)))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// wait until the inference is done
	c.JSON(200, gin.H{"message": "ok", "data": msg})
}
