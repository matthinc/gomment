package main

import (
	"os"
    "fmt"
    "syscall"

    "golang.org/x/crypto/ssh/terminal"
	
    "github.com/matthinc/gomment/api"
    "github.com/matthinc/gomment/persistence"
    "github.com/matthinc/gomment/logic"
    "github.com/matthinc/gomment/auth"
)

func main() {
    if len(os.Args) > 1 {
        command := os.Args[1]
        switch command {
        case "genpw":
            genpw()
        }
    }
    
	db := persistence.DB {}

    dbPath := os.Getenv("GOMMENT_DB_PATH")
    if len(dbPath) == 0 {
        dbPath = "./gomment.db"
    }

    pwHash := os.Getenv("GOMMENT_PW_HASH")
    if len(pwHash) == 0 {
        fmt.Println("admin password hash variable was not provided (GOMMENT_PW_HASH), disabling administration")
    }
    
	err := db.Open(dbPath)
	if err != nil {
		os.Exit(2)
	}
	err = db.Setup()
	if err != nil {
		os.Exit(3)
	}

    logic := logic.BusinessLogic { &db, pwHash }
    
    api.StartApi(&logic)
}

func genpw() {
    fmt.Print("Enter your password: ")
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        fmt.Println("failed to read password")
        os.Exit(1)
    }
    hash, err := auth.HashPw(string(bytePassword))
    if err != nil {
        fmt.Println("failed to hash password")
        os.Exit(1)
    }
    fmt.Println("\n" + hash)
    os.Exit(0)
}
