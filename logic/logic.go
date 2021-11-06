package logic

import (
	"time"

	"github.com/matthinc/gomment/persistence"
)

type SessionData struct {
	ValidUntil time.Time
}

type BusinessLogic struct {
	DB         persistence.Persistence
	PwHash     string
	SessionMap map[string]SessionData
}

func Create(db persistence.Persistence, pwHash string) BusinessLogic {
	return BusinessLogic{
		db,
		pwHash,
		make(map[string]SessionData),
	}
}
