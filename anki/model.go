package anki

func NewNote(front, back string, tags []string, example string) *Note {
	noteOptions := &noteOptions{
		AllowDuplicate: false,
		DuplicateScope: "deck",
		DuplicateScopeOptions: &duplicateScopeOptions{
			DeckName: deckName,
		},
	}

	return &Note{
		DeckName:  deckName,
		ModelName: modelName,
		Fields: &noteFields{
			Front:   front,
			Back:    back,
			Example: example,
		},
		Options: noteOptions,
		Tags:    tags,
	}
}

type NoteID int

type Note struct {
	NoteID    NoteID       `json:"noteId,omitempty"`
	DeckName  string       `json:"deckName"`
	ModelName string       `json:"modelName"`
	Fields    *noteFields  `json:"fields"`
	Options   *noteOptions `json:"options,omitempty"`
	Tags      []string     `json:"tags,omitempty"`
	Audio     []*noteMedia `json:"audio,omitempty"`
	Video     []*noteMedia `json:"video,omitempty"`
	Picture   []*noteMedia `json:"picture,omitempty"`

	Exists bool `json:"-"`
}

type noteOptions struct {
	AllowDuplicate        bool                   `json:"allowDuplicate"`
	DuplicateScope        string                 `json:"duplicateScope"`
	DuplicateScopeOptions *duplicateScopeOptions `json:"duplicateScopeOptions"`
}

type duplicateScopeOptions struct {
	DeckName string `json:"deckName"`
}

type noteFields struct {
	Front   string `json:"Front,omitempty"`
	Back    string `json:"Back,omitempty"`
	Example string `json:"Example,omitempty"`
}

type noteMedia struct {
	URL      string   `json:"url"`
	Filename string   `json:"filename"`
	SkipHash string   `json:"skipHash,omitempty"`
	Fields   []string `json:"fields"`
}

func NewNoteMedia(url, Filename string, fields ...string) *noteMedia {
	return &noteMedia{
		URL:      url,
		Filename: Filename,
		Fields:   fields,
	}
}

type Notes []*Note

func (notes Notes) GetByWord(word string) *Note {
	for _, n := range notes {
		if n.Fields.Front == word {
			return n
		}
	}

	return nil
}

func (notes Notes) FindByWords(words []string) Notes {
	result := Notes{}

	for _, w := range words {
		if note := notes.GetByWord(w); note != nil {
			result = append(result, note)
		}
	}

	return result
}

func (notes Notes) ListWords() []string {
	words := make([]string, len(notes))

	for i, n := range notes {
		words[i] = n.Fields.Front
	}

	return words
}

type action struct {
	Action string      `json:"action"`
	Params interface{} `json:"params,omitempty"`
}

func NewSyncAction() *action {
	return &action{Action: "sync"}
}

func NewAddNotesAction(notes Notes) *action {
	return &action{
		Action: "addNotes",
		Params: &addNotesParams{
			Notes: notes,
		},
	}
}

type findNotesParams struct {
	Query string `json:"query"`
}

type notesInfoParams struct {
	NoteIDs []NoteID `json:"notes"`
}

type addNotesParams struct {
	Notes Notes `json:"notes"`
}

type multiParams struct {
	Actions []*action `json:"actions"`
}

type canAddNotesParams struct {
	Notes Notes `json:"notes"`
}

type updateNoteFieldsParams struct {
	Note *struct {
		ID     NoteID      `json:"id"`
		Fields *noteFields `json:"fields"`
	} `json:"note"`
}

type params interface {
	findNotesParams | notesInfoParams | addNotesParams | canAddNotesParams | multiParams | updateNoteFieldsParams | struct{}
}

type noteInfoResult struct {
	NoteID    NoteID   `json:"noteId"`
	ModelName string   `json:"modelName"`
	Tags      []string `json:"tags"`
	Fields    struct {
		Front struct {
			Order int    `json:"order"`
			Value string `json:"value"`
		} `json:"Front"`
		Back struct {
			Order int    `json:"order"`
			Value string `json:"value"`
		} `json:"Back"`
		Example struct {
			Order int    `json:"order"`
			Value string `json:"value"`
		} `json:"Example"`
	} `json:"fields"`
	Cards []int `json:"cards"`
}
