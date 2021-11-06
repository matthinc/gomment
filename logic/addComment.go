package logic

import (
	"fmt"
	"html"
	"time"

	"github.com/matthinc/gomment/model"
)

func sanitize(commentCreation *model.CommentCreation) *model.CommentCreation {
	commentCreation.Author = html.EscapeString(commentCreation.Author)
	commentCreation.Text = html.EscapeString(commentCreation.Text)
	commentCreation.Email = html.EscapeString(commentCreation.Email)

	return commentCreation
}

func (logic *BusinessLogic) CreateComment(commentCreation *model.CommentCreation) (int64, error) {
	commentCreation = sanitize(commentCreation)

	id, err := logic.DB.CreateComment(commentCreation, time.Now().Unix())
	if err != nil {
		return 0, fmt.Errorf("unable to create comment in the database: %w", err)
	}

	return id, nil
}
