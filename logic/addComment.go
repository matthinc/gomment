package logic

import (
    "github.com/matthinc/gomment/model"
    "fmt"
)

func (logic* BusinessLogic) AddComment(comment *model.Comment) {
    err := logic.DB.AddComment(comment)
    if err != nil {
        fmt.Println(err)
    }
}
