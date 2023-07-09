package database

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// MustNotEmptyString panics if the string is empty.
func MustNotEmptyString(s string) string {
	if s == "" {
		panic("string must not be empty")
	}
	return s
}

// MustLoadEnv try to load .env file from a path.
// If file does not exist, it will do nothing.
func MustLoadEnv(path ...string) {
	if err := godotenv.Load(path...); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return
		}

		panic("Failed to load .env file")
	}
}

// SearchUpwardForFile searches for a file in the current directory and all parent directories.
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
