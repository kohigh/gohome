package main

import (
	"log"
	"os"
)

func main() {
	path := os.Args[1]

	envs, err := ReadDir(path)
	if err != nil {
		log.Fatalf("Problem with opening file: %v", err)
	}

	RunCmd(os.Args[2:], envs)
}
