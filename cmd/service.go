package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ezio1119/add-anki-flashcards-with-weblio/anki"
	"github.com/ezio1119/add-anki-flashcards-with-weblio/weblio"
	"golang.org/x/sync/errgroup"
)

func addWords(ctx context.Context, words []string) error {

	wg := errgroup.Group{}
	wg.SetLimit(10)

	notes := []*anki.Note{}

	for i, w := range words {
		i := i
		w := w

		wg.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, time.Microsecond*5)
			defer cancel()

			result, err := queryWordWeblio(ctx, w)
			if err != nil {
				return err
			}

			front := result.Query
			back := strings.TrimSpace(result.Description)
			tags := []string{}
			if result.Level != 0 {
				tags = append(tags, strconv.Itoa(result.Level))
			}

			note := anki.NewNote(front, back, tags)

			if result.AudioURL != "" {
				audio := &anki.NoteMedia{
					URL:      result.AudioURL,
					Filename: front,
					Fields:   []string{"Front"},
				}

				note.Audio = append(note.Audio, audio)
			}

			notes = append(notes, note)

			fmt.Printf("progress on weblio %d/%d\n", i+1, len(words))
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	newNotes, err := removeAnkiDupNotes(ctx, notes)
	if err != nil {
		return err
	}

	if err := addNotesAnki(ctx, newNotes); err != nil {
		return err
	}

	fmt.Printf("added: %d: duplicated: %d: received: %d \n\n", len(newNotes), len(words)-len(newNotes), len(words))

	for _, n := range notes {
		fmt.Printf("%s: %s\n", n.Fields.Front, n.Fields.Back)
	}

	return nil
}

func queryWordWeblio(ctx context.Context, word string) (*weblio.QueryResult, error) {
	result, err := weblio.Query(ctx, word)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func removeAnkiDupNotes(ctx context.Context, notes []*anki.Note) ([]*anki.Note, error) {
	if err := anki.CanAddNotes(ctx, notes); err != nil {
		return nil, err
	}

	newNotes := []*anki.Note{}

	for _, n := range notes {
		if n.CanAdd {

			wordWithAudio := fmt.Sprintf("%s[sound:%s]", n.Fields.Front, n.Fields.Front)
			noteWithAudio := anki.NewNote(wordWithAudio, "", nil)

			if err := anki.CanAddNotes(ctx, []*anki.Note{noteWithAudio}); err != nil {
				return nil, err
			}

			if noteWithAudio.CanAdd {
				newNotes = append(newNotes, n)
			}
		}

	}

	return newNotes, nil
}

func addNotesAnki(ctx context.Context, notes []*anki.Note) error {
	actions := make([]*anki.Action, 2)

	addNotesAction := &anki.Action{
		Action: "addNotes",
		Params: &anki.AddNotesParams{Notes: notes},
	}
	actions[0] = addNotesAction

	syncAction := &anki.Action{Action: "sync"}
	actions[1] = syncAction

	multiParams := &anki.MultiParams{Actions: actions}

	if err := anki.Multi(ctx, multiParams); err != nil {
		return err
	}

	return nil
}
