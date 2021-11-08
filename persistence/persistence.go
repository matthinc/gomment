package persistence

import (
	"github.com/matthinc/gomment/model"
)

type ThreadMetaInfo struct {
	NumTotal int
	NumRoot  int
	ThreadId int64
}

type Persistence interface {
	Open(path string) error
	Setup() error
	Close()
	CreateComment(commentCreation *model.CommentCreation, createdAt int64) (int64, error)
	GetNewestCommentsByPath(path string, limit int) ([]model.Comment, ThreadMetaInfo, error)
	GetMoreNewestSiblings(threadId int64, parentId int64, newestCreatedAt int64, excludeIds []int64, limit int) ([]model.Comment, error)
	GetThreads() ([]model.Thread, error)
}
