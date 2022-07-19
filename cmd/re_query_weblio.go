package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ezio1119/add-anki-flashcards-with-weblio/anki"
	"github.com/ezio1119/add-anki-flashcards-with-weblio/util"
	"github.com/ezio1119/add-anki-flashcards-with-weblio/weblio"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

func RunReQueryWeblio(ctx context.Context) error {
	noteIDs, err := anki.FindNotes(ctx)
	if err != nil {
		return err
	}

	existsNotes, err := anki.NotesInfo(ctx, noteIDs)
	if err != nil {
		return err
	}

	var notes anki.Notes
	for _, n := range existsNotes {
		if n.Fields.Example == "" {
			notes = append(notes, n)
		}
	}

	eg := errgroup.Group{}
	eg.SetLimit(5)

	bar := progressbar.Default(int64(len(notes)))

	for i := range notes {
		i := i

		eg.Go(func() error {

			w := notes[i].Fields.Front
			w = util.RemoveAudioFromWord(w)

			weblioCtx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

			r, err := weblio.Query(weblioCtx, w)
			if err != nil {
				fmt.Printf("[ERROR] %s: failed weblio query: %s\n", w, err)
				return err
			}

			notes[i].Fields.Front = r.Query
			notes[i].Fields.Back = r.Description
			notes[i].Fields.Example = r.Examples.String()

			if r.AudioURL != "" {
				notes[i].Audio = nil
				audio := anki.NewNoteMedia(r.AudioURL, r.Query, "Front")
				notes[i].Audio = append(notes[i].Audio, audio)
			}

			if err := anki.UpdateNoteFields(ctx, notes[i]); err != nil {
				fmt.Printf("[ERROR] %s: failed update note fields: %s\n", w, err)
				return err
			}

			if err := bar.Add(1); err != nil {
				fmt.Println(err)
			}

			return nil
		})
	}

	return eg.Wait()
}
