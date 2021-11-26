package api

import "github.com/matthinc/gomment/logic"

type ServerConfig struct {
	RequireEmail    bool `json:"require_email"`
	RequireAuthor   bool `json:"require_author"`
	RequireApproval bool `json:"require_approval"`
}

type ThreadCommentsResponse struct {
	Config ServerConfig        `json:"config"`
	Thread logic.CommentResult `json:"thread"`
}
