package persistence

import (
	"github.com/matthinc/gomment/model"
)

type Persistence interface {
	Open(path string) error
	Setup() error
	Close()
	AddComment(comment* model.Comment) (int64, error)
    QueryComments(thread int) []model.Comment
    QueryThreads() []model.Thread
}
