package logic

import (
	"fmt"
	"html"

	"github.com/matthinc/gomment/model"
)

func sanitize(commentCreation *model.CommentCreation) *model.CommentCreation {
	commentCreation.Author = html.EscapeString(commentCreation.Author)
	commentCreation.Text = html.EscapeString(commentCreation.Text)
	commentCreation.Email = html.EscapeString(commentCreation.Email)

	return commentCreation
}

func (logic *BusinessLogic) CreateComment(commentCreation *model.CommentCreation) int64 {
	commentCreation = sanitize(commentCreation)
	id, err := logic.DB.CreateComment(commentCreation)
	if err != nil {
		fmt.Println(err)
	}

	return id
}
