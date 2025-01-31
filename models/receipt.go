package models

import "github.com/google/uuid"

type Receipt struct {
	Id           uuid.UUID `json:"id"`
	Retailer     string    `json:"retailer" binding:"required"`
	PurchaseDate string    `json:"purchaseDate" binding:"required"`
	PurchaseTime string    `json:"purchaseTime" binding:"required"`
	Items        []Item    `json:"items" binding:"required"`
	Total        string    `json:"total" binding:"required"`
}
