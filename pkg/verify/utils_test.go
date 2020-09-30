package verify

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "studio-abuse-detector"

// LoadEnv loads env vars from .env
// https://github.com/joho/godotenv/issues/43
func LoadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Problem loading .env file", err, cwd)
		os.Exit(-1)
	}
}
