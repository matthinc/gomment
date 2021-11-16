package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/matthinc/gomment/api"
	"github.com/matthinc/gomment/auth"
	"github.com/matthinc/gomment/logic"
	"github.com/matthinc/gomment/persistence/sqlite"
)

func main() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := loggerConfig.Build()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "genpw":
			genpw()
		}
	}

	db := sqlite.DB{}

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
		log.Fatal(err)
		os.Exit(2)
	}
	err = db.Setup()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	logic := logic.Create(&db, pwHash)

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
