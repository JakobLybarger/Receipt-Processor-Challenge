package main

import (
	"github.com/JakobLybarger/ReceiptProcessorChallenge/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/receipts/process", handlers.ProcessReceipt)
	r.GET("/receipts/:id/points", handlers.CalculatePoints)

	r.Run()
}
