package cmd

import (
	"context"
	"log"
	"os"

	"github.com/ezio1119/add-anki-flashcards-with-weblio/util"
)

func RunCLI(ctx context.Context) error {
	args := os.Args

	for len(args) < 2 {
		log.Fatal("argument words required")
	}

	words := util.RemoveElem(args, 0)

	return addWords(ctx, words)
}
