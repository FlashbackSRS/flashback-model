package fb

// Package represents a top-level collection of Flashback Documents, such that
// they can be easily transmitted or shared as a single file. It is intended to
// be used via its json.Marshaler and json.Unmarshaler interfaces.
type Package struct {
	Bundle  *Bundle   `json:"bundle,omitempty"`
	Cards   []*Card   `json:"cards,omitempty"`
	Notes   []*Note   `json:"notes,omitempty"`
	Decks   []*Deck   `json:"decks,omitempty"`
	Themes  []*Theme  `json:"themes,omitempty"`
	Reviews []*Review `json:"reviews,omitempty"`
}
