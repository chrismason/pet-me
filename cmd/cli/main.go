package main

import (
	"log"

	"github.com/chrismason/pet-me/cmd/cli/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
