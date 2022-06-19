package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ezio1119/add-anki-flashcards-with-weblio/anki"
	"github.com/ezio1119/add-anki-flashcards-with-weblio/weblio"
	"golang.org/x/sync/errgroup"
)

func addWords(ctx context.Context, words []string) error {

	notes := make([]*anki.Note, len(words))

	for i, w := range words {

		note := anki.NewNote(w, "", nil, nil, nil, nil)
		notes[i] = note
	}

	if err := anki.CanAddNotes(ctx, notes); err != nil {
		return err
	}

	canAddNotes := []*anki.Note{}
	for _, n := range notes {
		if n.CanAdd {
			canAddNotes = append(canAddNotes, n)
		}
	}

	wg := errgroup.Group{}
	wg.SetLimit(10)

	filledNotes := []*anki.Note{}

	for i, n := range canAddNotes {
		i := i
		n := n
		wg.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, time.Microsecond*5)
			defer cancel()

			result, err := queryWordWeblio(ctx, n.Fields.Front)
			if err != nil {
				return err
			}

			n.Fields.Front = result.Query
			n.Fields.Back = result.Description

			level := strconv.Itoa(result.Level)
			n.Tags = append(n.Tags, level)

			if result.AudioURL != "" {
				audio := &anki.NoteMedia{
					URL:      result.AudioURL,
					Filename: n.Fields.Front,
					Fields:   []string{"Front"},
				}

				n.Audio = append(n.Audio, audio)
			}

			filledNotes = append(filledNotes, n)

			fmt.Printf("weblio %d/%d\n", i+1, len(canAddNotes))

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	if err := addNotesAnki(ctx, filledNotes); err != nil {
		return err
	}

	fmt.Printf("success: %d: canAdd: %d: received: %d\n", len(filledNotes), len(canAddNotes), len(words))

	return nil
}

func queryWordWeblio(ctx context.Context, word string) (*weblio.QueryResult, error) {
	result, err := weblio.Query(ctx, word)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%#v\n", result)
	return result, nil
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
