package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ezio1119/add-anki-flashcards-with-weblio/anki"
	"github.com/ezio1119/add-anki-flashcards-with-weblio/util"
	"github.com/ezio1119/add-anki-flashcards-with-weblio/weblio"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"golang.org/x/sync/errgroup"
)

func addWords(ctx context.Context, words []string) error {
	for i, n := range words {
		words[i] = strings.ToLower(n)
	}

	existsNotes, err := findExistsNotesFromWords(ctx, words)
	if err != nil {
		return err
	}

	existsWords := existsNotes.ListWords()
	notExistsWords := make([]string, 0, len(words)-len(existsWords))

	for _, w := range words {
		var exists bool

		for _, wExists := range existsWords {
			if wExists == w {
				exists = true
			}
		}

		if !exists {
			notExistsWords = append(notExistsWords, w)
		}
	}

	wg := errgroup.Group{}
	wg.SetLimit(10)

	newNotes := []*anki.Note{}
	failedQueryWords := []string{}

	for i, w := range notExistsWords {
		i := i
		w := w

		wg.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, time.Microsecond*5)
			defer cancel()

			result, err := queryWordWeblio(ctx, w)
			if err != nil {
				fmt.Printf("addWords: failed query %s\n", w)
				failedQueryWords = append(failedQueryWords, w)
				return nil
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

			newNotes = append(newNotes, note)

			fmt.Printf("querying weblio... %d/%d\n", i+1, len(notExistsWords))

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	if len(newNotes) != 0 {
		if err := addNotesAnki(ctx, newNotes); err != nil {
			return err
		}
	}

	fmt.Printf("added: %d: duplicated: %d: failedQuery: %d: received: %d \n\n", len(newNotes), len(existsNotes), len(failedQueryWords), len(words))

	outputNotes(existsNotes)
	outputNotes(newNotes)

	return nil
}

func findExistsNotesFromWords(ctx context.Context, words []string) (anki.Notes, error) {
	noteIDs, err := anki.FindNotes(ctx)
	if err != nil {
		return nil, err
	}

	allNotes, err := anki.NotesInfo(ctx, noteIDs)
	if err != nil {
		return nil, err
	}

	// ankiに登録されてるのは出力し、されてないものはweblioに投げる
	existsNotes := allNotes.FindByWords(words)

	for _, w := range words {
		wordWithAudio := util.AddAudioToWord(w)

		existsNote := allNotes.GetByWord(wordWithAudio)

		if existsNote != nil {
			existsNote.Fields.Front = w
			existsNotes = append(existsNotes, existsNote)
		}
	}

	return existsNotes, nil
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
		if !n.Exists {

			wordWithAudio := util.AddAudioToWord(n.Fields.Front)
			noteWithAudio := anki.NewNote(wordWithAudio, "", nil)

			if err := anki.CanAddNotes(ctx, []*anki.Note{noteWithAudio}); err != nil {
				return nil, err
			}

			if !noteWithAudio.Exists {
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

func outputNotes(notes anki.Notes) {
	for _, n := range notes {
		fmt.Printf("%s: %s\n", n.Fields.Front, n.Fields.Back)

		// if len(n.Audio) == 1 {
		// 	if err := playSoundFromURL(n.Audio[0].URL); err != nil {
		// 		fmt.Printf("addWords: playSoundFromURL: failed tp play '%s' sound: %s\n", n.Fields.Front, err)
		// 	}
		// }
	}
}

func playSoundFromURL(url string) error {
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	d, err := mp3.NewDecoder(res.Body)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8182)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}

	time.Sleep(time.Second)

	return nil
}
