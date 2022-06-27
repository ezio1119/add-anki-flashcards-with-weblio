package anki

func NewNote(front, back string, tags []string) *Note {
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
			Front: front,
			Back:  back,
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
	Audio     []*NoteMedia `json:"audio,omitempty"`
	Video     []*NoteMedia `json:"video,omitempty"`
	Picture   []*NoteMedia `json:"picture,omitempty"`

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
	Front string `json:"Front"`
	Back  string `json:"Back"`
}

type NoteMedia struct {
	URL      string   `json:"url"`
	Filename string   `json:"filename"`
	SkipHash string   `json:"skipHash"`
	Fields   []string `json:"fields"`
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

type Action struct {
	Action string      `json:"action"`
	Params interface{} `json:"params,omitempty"`
}

type FindNotesParams struct {
	Query string `json:"query"`
}

type NotesInfoParams struct {
	NoteIDs []NoteID `json:"notes"`
}

type AddNotesParams struct {
	Notes Notes `json:"notes"`
}

type MultiParams struct {
	Actions []*Action `json:"actions"`
}

type CanAddNotesParams struct {
	Notes Notes `json:"notes"`
}

type params interface {
	FindNotesParams | NotesInfoParams | AddNotesParams | CanAddNotesParams | MultiParams | struct{}
}

type NoteInfoResult struct {
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
	} `json:"fields"`
	Cards []int `json:"cards"`
}
