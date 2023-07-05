package database

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func MustNotEmptyString(s string) string {
	if s == "" {
		log.Panic().Msg("string must not be empty")
	}
	return s
}

func MustLoadEnv(path ...string) {
	if err := godotenv.Load(path...); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Warn().Err(err).Msg("no .env file found")
			return
		}

		log.Fatal().Err(err).Msg("Failed to load .env file")
	}
}

func SearchUpwardForFile(fileName string) string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		absFilePath := filepath.Join(dir, fileName)
		if _, err := os.Stat(absFilePath); !os.IsNotExist(err) {
			return absFilePath
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return ""
		}

		dir = parentDir
	}
}
