package model

type Comment struct {
	Id          int64  `json:"comment_id"`
	Author      string `json:"author" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Text        string `json:"text" binding:"required"`
	ParentId    int64  `json:"parent_id" binding:""`
	CreatedAt   int64  `json:"created_at"`
	TouchedAt   int64  `json:"touched_at"`
	NumChildren int    `json:"num_children"`
}

type CommentCreation struct {
	ThreadPath string `json:"thread_path" binding:"required"`
	Author     string `json:"author"`
	Email      string `json:"email"`
	Text       string `json:"text"`
	ParentId   int64  `json:"parent_id"` // optional, maps to 0 if not provided
}

type Thread struct {
	Id       int64
	Path     string
	NumTotal int
	NumRoot  int
}
