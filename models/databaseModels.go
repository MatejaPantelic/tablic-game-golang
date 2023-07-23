package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `json:"name"`
}

type Game struct {
	gorm.Model
	Score  int  `json:"score"`
	UserID int  `json:"user_id"`
	DeckId string
	HandPile string
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
