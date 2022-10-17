package model

type Config struct {
	Model
	Content string `json:"content" binding:"required"`
}
