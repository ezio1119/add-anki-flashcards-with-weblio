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
		Fields: noteFields{
			Front: front,
			Back:  back,
		},
		Options: noteOptions,
		Tags:    tags,
	}
}

type Note struct {
	DeckName  string       `json:"deckName"`
	ModelName string       `json:"modelName"`
	Fields    noteFields   `json:"fields"`
	Options   *noteOptions `json:"options,omitempty"`
	Tags      []string     `json:"tags,omitempty"`
	Audio     []*NoteMedia `json:"audio,omitempty"`
	Video     []*NoteMedia `json:"video,omitempty"`
	Picture   []*NoteMedia `json:"picture,omitempty"`

	CanAdd bool `json:"-"`
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

type Action struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}

type AddNotesParams struct {
	Notes []*Note `json:"notes"`
}

type MultiParams struct {
	Actions []*Action `json:"actions"`
}

type CanAddNotesParams struct {
	Notes []*Note `json:"notes"`
}

type params interface {
	AddNotesParams | CanAddNotesParams | MultiParams | struct{}
}
