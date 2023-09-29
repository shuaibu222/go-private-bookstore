package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

var projectDirName = "go-bookstore"

func LoadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`) // Create a regular expression to match the project directory name
	currentWorkDirectory, _ := os.Getwd()                            // Get the current working directory and assign it to the currentWorkDirectory variable
	rootPath := projectName.Find([]byte(currentWorkDirectory))       // Use the regular expression to find the project directory in the current working directory
	err := godotenv.Load(string(rootPath) + `/.env`)                 // Load the .env file from the rootPath

	if err != nil {
		log.Fatalf("Error loading .env file") // Log a fatal error message and exit the program
	}
}
