package handlers

import (
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
			points := calculateReceiptPoints(receipt)
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

func calculateReceiptPoints(receipt models.Receipt) int {
	points := 0
	for _, r := range receipt.Retailer {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			points++
		}
	}

	fmt.Printf("Points after retiler name %d\n", points)

	prev := points

	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Trunc(total) {
		points += 50
	}

	fmt.Printf("Points after seeing if total is round dollar %d - %d \n", points, points-prev)
	prev = points

	if math.Mod(total, 0.25) == 0 {
		points += 25
	}
	fmt.Printf("Points after seeing if is multiple of 0.25 %d - %d\n", points, points-prev)
	prev = points

	for i := 1; i < len(receipt.Items); i += 2 {
		points += 5
	}
	fmt.Printf("Points after for every 2 items %d - %d\n", points, points-prev)
	prev = points

	for _, item := range receipt.Items {
		strLen := float64(len(strings.TrimSpace(item.ShortDescription)))
		if math.Mod(strLen, 3) == 0 {

			fmt.Printf("'%s'\n'%s'\n", item.ShortDescription, strings.TrimSpace(item.ShortDescription))
			price, _ := strconv.ParseFloat(item.Price, 64)
			fmt.Printf("price %f\n", price)
			newpoints := int(math.Ceil(price * 0.2))
			fmt.Printf("new points %d\n", newpoints)
			points += int(math.Ceil(price * 0.2))
		}
	}
	fmt.Printf("Points after trimmed description %d - %d\n", points, points-prev)
	prev = points

	d, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if math.Mod(float64(d.Day()), 2) != 0 {
		fmt.Printf("Days is %d", d.YearDay())
		points += 6
	}
	fmt.Printf("Points after checking if purchase date is odd %d - %d\n", points, points-prev)
	prev = points

	d, _ = time.Parse("15:04", receipt.PurchaseTime)
	fmt.Printf("Hour is %d\n", d.Hour())
	if 14 <= d.Hour() && d.Hour() < 16 {
		points += 10
	}
	fmt.Printf("Points after looking at purchase time %d - %d\n", points, points-prev)

	return points
}
