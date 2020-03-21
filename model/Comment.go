package model

type Comment struct {
    Name string    `json:"name" binding:"required"`
    Email string   `json:"email" binding:"required"`
    Content string `json:"content" binding:"required"`
}
