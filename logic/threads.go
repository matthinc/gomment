package logic

import (
	"fmt"

	"github.com/matthinc/gomment/model"
)

func (logic *BusinessLogic) GetThreads(thread int) ([]model.Thread, error) {
	threads, err := logic.DB.GetThreads()
	if err != nil {
		return []model.Thread{}, fmt.Errorf("unable to get list of threads from the database: %w", err)
	}
	return threads, nil
}
