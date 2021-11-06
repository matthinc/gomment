package persistence

import (
	"github.com/matthinc/gomment/model"
)

type Persistence interface {
	Open(path string) error
	Setup() error
	Close()
	CreateComment(commentCreation *model.CommentCreation) (int64, error)
	GetNewestCommentsByPath(path string, limit int) ([]model.Comment, error)
	GetThreads() ([]model.Thread, error)
}
