package logic

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/matthinc/gomment/model"
)

func sanitize(commentCreation *model.CommentCreation) *model.CommentCreation {
	commentCreation.Author = html.EscapeString(commentCreation.Author)
	commentCreation.Text = html.EscapeString(commentCreation.Text)
	commentCreation.Email = html.EscapeString(commentCreation.Email)

	return commentCreation
}

func ValidateComment(v ValidationT, commentCreation *model.CommentCreation) error {
	authorLength := uint(len(commentCreation.Author))
	if v.RequireAuthor && authorLength == 0 {
		return ValidationErrorRequired("author")
	}
	if authorLength != 0 && (authorLength < v.AuthorLengthMin || authorLength > v.AuthorLengthMax) {
		return ValidationErrorLength("author", v.AuthorLengthMin, v.AuthorLengthMax)
	}

	emailLength := uint(len(commentCreation.Email))
	if v.RequireEmail && emailLength == 0 {
		return ValidationErrorRequired("email")
	}
	if emailLength != 0 && (emailLength < v.EmailLengthMin || emailLength > v.EmailLengthMax) {
		return ValidationErrorLength("email", v.EmailLengthMin, v.EmailLengthMax)
	}
	if v.RequireEmail && !strings.Contains(commentCreation.Email, "@") {
		return ValidationErrorSymbol("email", "@")
	}

	textLength := len(commentCreation.Text)
	if textLength < int(v.CommentLengthMin) || textLength > int(v.CommentLengthMax) {
		return ValidationErrorLength("text", v.CommentLengthMin, v.CommentLengthMax)
	}

	return nil
}

func (logic *BusinessLogic) CreateComment(commentCreation *model.CommentCreation) (int64, error) {
	err := ValidateComment(logic.validation, commentCreation)
	if err != nil {
		return 0, err
	}

	commentCreation = sanitize(commentCreation)

	id, err := logic.DB.CreateComment(commentCreation, time.Now().Unix())
	if err != nil {
		return 0, fmt.Errorf("unable to create comment in the database: %w", err)
	}

	return id, nil
}
