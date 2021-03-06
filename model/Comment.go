package model

type Comment struct {
    Id int           `json:"comment_id"`
    Author string    `json:"author" binding:"required"`
    Email string     `json:"email" binding:"required"`
    Text string      `json:"text" binding:"required"`
    ThreadId int     `json:"thread_id" binding:""`
    ParentId int     `json:"parent_id" binding:""`
    Created string   `json:"created_at"`
}

type CommentTree struct {
    Children []CommentTree   `json:"children"`
    Comment Comment          `json:"comment"`
    HasChildren bool         `json:"has_children"`
}

type CommentsResponse struct {
    Comments []CommentTree   `json:"comments"`
    Total int                `json:"total"`
}

type Thread struct {
    Id int
}
