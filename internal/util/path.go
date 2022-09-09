package util

import (
	"log"
	"os"
	"path/filepath"
)

func RootDir() string {
	root, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get root dir, err: %s", err)
	}
	return filepath.Dir(root)
}
