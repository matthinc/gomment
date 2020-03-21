package main

import (
	"os"
	
    "github.com/matthinc/gomment/api"
    "github.com/matthinc/gomment/persistence"	
)

func main() {
	db := persistence.DB {}
	err := db.Open("./gomment.db")
	if err != nil {
		os.Exit(2)
	}
	err = db.Setup()
	if err != nil {
		os.Exit(3)
	}
	
    api.StartApi()
}
