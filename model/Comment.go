package model

type Comment struct {
	Id          int    `json:"comment_id"`
	Author      string `json:"author" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Text        string `json:"text" binding:"required"`
	ParentId    int    `json:"parent_id" binding:""`
	CreatedAt   int64  `json:"created_at"`
	TouchedAt   int64  `json:"touched_at"`
	NumChildren int    `json:"num_children"`
}

type CommentCreation struct {
	ThreadPath string `json:"thread_path" binding:"required"`
	Author     string `json:"author" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Text       string `json:"text" binding:"required"`
	ParentId   int    `json:"parent_id"` // optional, maps to 0 if not provided
}

type CommentTree struct {
	Children []CommentTree `json:"children"`
	Comment  Comment       `json:"comment"`
}

type CommentsResponse struct {
	Comments       []CommentTree `json:"comments"`
	NumRoot        int           `json:"num_root"`
	NumTotal       int           `json:"num_total"`
	NumRootPayload int           `json:"num_root_payload"`
	ThreadId       int64         `json:"thread_id"`
}

type Thread struct {
	Id       int
	Path     string
	NumTotal int
	NumRoot  int
}
