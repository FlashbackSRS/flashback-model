package fb

type Package struct {
	Bundle  *Bundle   `json:"bundle,omitempty"`
	Cards   []*Card   `json:"cards,omitempty"`
	Notes   []*Note   `json:"notes,omitempty"`
	Decks   []*Deck   `json:"decks,omitempty"`
	Themes  []*Theme  `json:"themes,omitempty"`
	Reviews []*Review `json:"reviews,omitempty"`
}
