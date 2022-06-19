package cmd

import (
	"context"
	"encoding/csv"
	"io"
	"os"
)

func RunImportCSV(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	r := csv.NewReader(f)

	var words []string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		word := record[3]
		if word == "" {
			continue
		}

		words = append(words, word)
	}

	// fmt.Printf("%#v\n", words)

	return addWords(ctx, words)
	// return nil
}
