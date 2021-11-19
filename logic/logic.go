package logic

import (
	"time"

	"github.com/matthinc/gomment/persistence"
)

type SessionData struct {
	ValidUntil time.Time
}

type AdministrationT struct {
	PasswordHash string
}

type ValidationT struct {
	RequireAuthor        bool
	RequireEmail         bool
	AuthorLengthMin      uint
	AuthorLengthMax      uint
	EmailLengthMin       uint
	EmailLengthMax       uint
	CommentLengthMin     uint
	CommentLengthMax     uint
	CommentDepthMax      uint
	InitialQueryDepthMax uint
	QueryLimitMax        uint
}

type BusinessLogic struct {
	DB             persistence.Persistence
	SessionMap     map[string]SessionData
	Administration AdministrationT
	validation     ValidationT
}

func GetDefaultValidation() ValidationT {
	return ValidationT{
		RequireAuthor:        false,
		RequireEmail:         false,
		AuthorLengthMin:      1,
		AuthorLengthMax:      50,
		EmailLengthMin:       3,
		EmailLengthMax:       100,
		CommentLengthMin:     1,
		CommentLengthMax:     20000,
		CommentDepthMax:      9,
		InitialQueryDepthMax: 3,
		QueryLimitMax:        200,
	}
}

func Create(db persistence.Persistence, administration AdministrationT, validation ValidationT) BusinessLogic {
	return BusinessLogic{
		db,
		make(map[string]SessionData),
		administration,
		validation,
	}
}
