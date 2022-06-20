package anki

import (
	"context"
	"encoding/json"
	"fmt"
)

func FindNotes(ctx context.Context) ([]NoteID, error) {
	query := fmt.Sprintf("deck:%s", deckName)

	res, err := requestAPI(ctx, "findNotes", &FindNotesParams{Query: query})
	if err != nil {
		fmt.Println("aa")
		return nil, err
	}

	noteIDs := []NoteID{}

	if err := json.Unmarshal(res.Result, &noteIDs); err != nil {
		return nil, err
	}

	return noteIDs, nil
}

func NotesInfo(ctx context.Context, noteIDs []NoteID) (Notes, error) {
	res, err := requestAPI(ctx, "notesInfo", &NotesInfoParams{NoteIDs: noteIDs})
	if err != nil {
		return nil, err
	}

	notesInfo := []*NoteInfoResult{}
	if err := json.Unmarshal(res.Result, &notesInfo); err != nil {
		return nil, err
	}

	// result := res.Result.()
	notes := make(Notes, len(notesInfo))

	for i, n := range notesInfo {

		note := NewNote(n.Fields.Front.Value, n.Fields.Back.Value, n.Tags)
		note.Exists = true
		notes[i] = note
	}

	return notes, nil
}

func AddNotes(ctx context.Context, notes Notes) error {
	_, err := requestAPI(ctx, "addNotes", &AddNotesParams{notes})
	return err
}

func Sync(ctx context.Context) error {
	_, err := requestAPI[struct{}](ctx, "sync", nil)
	return err
}

func Multi(ctx context.Context, params *MultiParams) error {
	_, err := requestAPI(ctx, "multi", params)
	return err
}

func CanAddNotes(ctx context.Context, notes Notes) error {
	res, err := requestAPI(ctx, "canAddNotes", &CanAddNotesParams{notes})
	if err != nil {
		return err
	}

	canAdds := []bool{}

	if err := json.Unmarshal(res.Result, &canAdds); err != nil {
		return err
	}

	if len(notes) != len(canAdds) {
		return fmt.Errorf("ERROR: anki: CanAddNotes: notes and canAdds are not the same length")
	}

	for i, can := range canAdds {
		if !can {
			notes[i].Exists = true
		}
	}

	return nil
}
