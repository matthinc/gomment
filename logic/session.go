package logic

import (
	"errors"
	"time"

	"github.com/matthinc/gomment/util"
)

const SessionDuration = time.Hour * time.Duration(1)

func (logic *BusinessLogic) CreateSession() (id string, data SessionData, err error) {
	id, err = util.GenerateRandomBase64(32)
	if err != nil {
		return "", SessionData{}, err
	}
	data = SessionData{
		ValidUntil: time.Now().Add(SessionDuration),
	}

	logic.SessionMap[id] = data

	return id, data, nil
}

func (logic *BusinessLogic) GetSession(id string) (data SessionData, err error) {
	d, ok := logic.SessionMap[id]
	if !ok {
		return SessionData{}, errors.New("No session exists for id")
	}
	if time.Now().Before(d.ValidUntil) {
		return d, nil
	} else {
		delete(logic.SessionMap, id)
		return SessionData{}, errors.New("Session has expired")
	}
}
