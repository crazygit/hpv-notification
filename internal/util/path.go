package util

import (
	"log"
	"os"
)

func RootDir() string {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get root dir, err: %s", err)
	}
	return root
}
