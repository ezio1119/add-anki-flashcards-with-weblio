package anki

func NewNote(front, back string, tags []string, audio, video, picture []*NoteMedia) *Note {
	return &Note{
		DeckName:  deckName,
		ModelName: modelName,
		Fields: noteFields{
			Front: front,
			Back:  back,
		},
		Tags:    tags,
		Audio:   audio,
		Video:   video,
		Picture: picture,
	}
}

type Note struct {
	DeckName  string       `json:"deckName"`
	ModelName string       `json:"modelName"`
	Fields    noteFields   `json:"fields"`
	Tags      []string     `json:"tags,omitempty"`
	Audio     []*NoteMedia `json:"audio,omitempty"`
	Video     []*NoteMedia `json:"video,omitempty"`
	Picture   []*NoteMedia `json:"picture,omitempty"`

	CanAdd, Added bool
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
