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

	noteIDs, err := anki.FindNotes(ctx)
	if err != nil {
		return err
	}

	allNotes, err := anki.NotesInfo(ctx, noteIDs)
	if err != nil {
		return err
	}

	existsNotes, err := findExistsNotesFromWords(ctx, allNotes, words)
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

	eg := errgroup.Group{}
	eg.SetLimit(10)

	newNotes := []*anki.Note{}
	failedQueryWords := []string{}

	for i, w := range notExistsWords {
		i := i
		w := w

		eg.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

			note, err := getNoteFromWeblio(ctx, w)
			if err != nil {
				fmt.Println(err)
				failedQueryWords = append(failedQueryWords, w)
				return err
			}

			newNotes = append(newNotes, note)

			fmt.Printf("querying weblio... %d/%d\n", i+1, len(notExistsWords))
			return nil
		})
	}

	eg.Wait()

	if len(newNotes) != 0 {
		if err := anki.Multi(ctx, anki.NewAddNotesAction(newNotes), anki.NewSyncAction()); err != nil {
			return err
		}
	}

	fmt.Printf("added: %d: duplicated: %d: failedQuery: %d: received: %d \n\n", len(newNotes), len(existsNotes), len(failedQueryWords), len(words))

	outputNotes(existsNotes)
	outputNotes(newNotes)

	return nil
}

func findExistsNotesFromWords(ctx context.Context, notes anki.Notes, words []string) (anki.Notes, error) {
	existsNotes := notes.FindByWords(words)

	for _, w := range words {
		wordWithAudio := util.AddAudioToWord(w)

		existsNote := notes.GetByWord(wordWithAudio)

		if existsNote != nil {
			existsNote.Fields.Front = w
			existsNotes = append(existsNotes, existsNote)
		}
	}

	return existsNotes, nil
}

func getNoteFromWeblio(ctx context.Context, w string) (*anki.Note, error) {
	result, err := weblio.Query(ctx, w)
	if err != nil {
		return nil, err
	}

	front := result.Query
	var tags []string
	if result.Level != 0 {
		tags = append(tags, strconv.Itoa(result.Level))
	}

	note := anki.NewNote(front, result.Description, tags, result.Examples.String())

	if result.AudioURL != "" {
		audio := anki.NewNoteMedia(result.AudioURL, front, "Front")
		note.Audio = append(note.Audio, audio)
	}

	return note, nil
}

func removeAnkiDupNotes(ctx context.Context, notes []*anki.Note) ([]*anki.Note, error) {
	if err := anki.CanAddNotes(ctx, notes); err != nil {
		return nil, err
	}

	newNotes := []*anki.Note{}

	for _, n := range notes {
		if !n.Exists {

			wordWithAudio := util.AddAudioToWord(n.Fields.Front)
			noteWithAudio := anki.NewNote(wordWithAudio, "", nil, "")

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
