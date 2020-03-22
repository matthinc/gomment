package main

import (
	"os"
	
    "github.com/matthinc/gomment/api"
    "github.com/matthinc/gomment/persistence"
    "github.com/matthinc/gomment/logic"	
)

func main() {
	db := persistence.DB {}

    dbPath := os.Getenv("GOMMENT_DB_PATH")
    if len(dbPath) == 0 {
        dbPath = "./gomment.db"
    }
    
	err := db.Open(dbPath)
	if err != nil {
		os.Exit(2)
	}
	err = db.Setup()
	if err != nil {
		os.Exit(3)
	}

    logic := logic.BusinessLogic { &db }
    
    api.StartApi(&logic)
}
