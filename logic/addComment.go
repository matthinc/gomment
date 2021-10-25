package logic

import (
	"fmt"
	"html"

	"github.com/matthinc/gomment/model"
)

func sanitize(comment *model.Comment) *model.Comment {
	comment.Author = html.EscapeString(comment.Author)
	comment.Text = html.EscapeString(comment.Text)
	comment.Email = html.EscapeString(comment.Email)

	return comment
}

func (logic *BusinessLogic) AddComment(comment *model.Comment) int64 {
	comment = sanitize(comment)
	id, err := logic.DB.AddComment(comment)
	if err != nil {
		fmt.Println(err)
	}

	return id
}
