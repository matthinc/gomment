package persistence

import (
	"github.com/matthinc/gomment/model"
)

type ThreadMetaInfo struct {
	NumTotal int
	NumRoot  int
}

type Persistence interface {
	Open(path string) error
	Setup() error
	Close()
	CreateComment(commentCreation *model.CommentCreation, createdAt int64) (int64, error)
	GetNewestCommentsByPath(path string, limit int) ([]model.Comment, ThreadMetaInfo, error)
	GetThreads() ([]model.Thread, error)
}
