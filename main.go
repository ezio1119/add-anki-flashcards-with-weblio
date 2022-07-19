package main

import (
	"context"
	"log"

	"github.com/ezio1119/add-anki-flashcards-with-weblio/cmd"
)

func main() {
	ctx := context.Background()
	// if err := cmd.RunImportCSV(ctx, "words.csv"); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := cmd.RunReQueryWeblio(ctx); err != nil {
	// 	log.Fatal(err)
	// }
	if err := cmd.RunCLI(ctx); err != nil {
		log.Fatal(err)
	}
}
