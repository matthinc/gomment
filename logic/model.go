package logic

import "github.com/matthinc/gomment/model"

type CommentTree struct {
	Children []CommentTree `json:"children"`
	Comment  model.Comment `json:"comment"`
}

type CommentResult struct {
	Comments       []CommentTree `json:"comments"`
	NumRoot        int           `json:"num_root"`
	NumTotal       int           `json:"num_total"`
	NumRootPayload int           `json:"num_root_payload"`
	ThreadId       int64         `json:"thread_id"`
}
