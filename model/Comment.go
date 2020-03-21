package model

type Comment struct {
    Author string    `json:"author" binding:"required"`
    Email string     `json:"email" binding:"required"`
    Text string      `json:"text" binding:"required"`
    ThreadId int     `json:"thread_id" binding:""`
    ParentId int     `json:"parent_id" binding:""`
}
