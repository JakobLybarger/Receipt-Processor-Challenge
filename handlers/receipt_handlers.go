package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/JakobLybarger/ReceiptProcessorChallenge/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var receipts []models.Receipt

func ProcessReceipt(c *gin.Context) {

	var receipt models.Receipt
	if err := c.ShouldBind(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "The receipt is invalid.",
		})
		return
	}

	receipt.Id = uuid.New()
	receipts = append(receipts, receipt)

	c.JSON(http.StatusCreated, gin.H{
		"id": receipt.Id,
	})
}

func CalculatePoints(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "Invalid ID",
		})

		return
	}

	for _, receipt := range receipts {
		if receipt.Id == id {
			points, err := calculateReceiptPoints(receipt)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"description": fmt.Sprintf("Unfortunately, there was bad data in the body... Error: %v", err.Error()),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"points": points,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"description": "No receipt found for that ID.",
	})
}

func calculateReceiptPoints(receipt models.Receipt) (int, error) {
	points := 0
	for _, r := range receipt.Retailer {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			points++
		}
	}

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return -1, err
	}

	if total == math.Trunc(total) {
		points += 50
	}

	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	for i := 1; i < len(receipt.Items); i += 2 {
		points += 5
	}

	for _, item := range receipt.Items {
		strLen := float64(len(strings.TrimSpace(item.ShortDescription)))
		if math.Mod(strLen, 3) == 0 {

			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	d, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return -1, err
	}

	if math.Mod(float64(d.Day()), 2) != 0 {
		points += 6
	}

	d, err = time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return -1, err
	}

	if 14 <= d.Hour() && d.Hour() < 16 {
		points += 10
	}

	return points, nil
}
