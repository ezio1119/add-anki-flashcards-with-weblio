package anki

import (
	"context"
	"fmt"
)

func AddNotes(ctx context.Context, notes []*Note) error {
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

func CanAddNotes(ctx context.Context, notes []*Note) error {
	res, err := requestAPI(ctx, "canAddNotes", &CanAddNotesParams{notes})
	if err != nil {
		return err
	}

	results := res.Result.([]interface{})
	canAdds := make([]bool, len(results))

	for i, r := range results {
		canAdds[i] = r.(bool)
	}

	if len(notes) != len(canAdds) {
		return fmt.Errorf("ERROR: anki: CanAddNotes: notes and canAdds are not the same length")
	}

	for i, can := range canAdds {
		if can {
			notes[i].CanAdd = can
		}
	}

	return nil
}
