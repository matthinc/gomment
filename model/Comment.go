package model

type Comment struct {
	Id          int    `json:"comment_id"`
	Author      string `json:"author" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Text        string `json:"text" binding:"required"`
	ThreadId    int    `json:"thread_id" binding:""`
	ParentId    int    `json:"parent_id" binding:""`
	CreatedAt   int64  `json:"created_at"`
	TouchedAt   int64  `json:"touched_at"`
	NumChildren int    `json:"num_children"`
}

type CommentCreation struct {
	ThreadPath string `json:"thread_path" binding:"required"`
	ParentId   int    `json:"parent_id" binding:"required"`
	Author     string `json:"author" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Text       string `json:"text" binding:"required"`
}

type CommentTree struct {
	Children []CommentTree `json:"children"`
	Comment  Comment       `json:"comment"`
}

type CommentsResponse struct {
	Comments    []CommentTree `json:"comments"`
	NumChildren int           `json:"num_children"`
	Total       int           `json:"total"`
}

type Thread struct {
	Id   int
	Path string
}
