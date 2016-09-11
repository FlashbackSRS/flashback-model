package fb

import "encoding/json"

type version int

const (
	// CurrentVersion represents the package format.
	CurrentVersion = 0
	// LowestVersion is the lowest version we can compatibly read.
	LowestVersion = 0
)

// Package represents a top-level collection of Flashback Documents, such that
// they can be easily transmitted or shared as a single file. It is intended to
// be used via its json.Marshaler and json.Unmarshaler interfaces.
type Package struct {
	Version int       `json:"version"`
	Bundle  *Bundle   `json:"bundle,omitempty"`
	Cards   []*Card   `json:"cards,omitempty"`
	Notes   []*Note   `json:"notes,omitempty"`
	Decks   []*Deck   `json:"decks,omitempty"`
	Themes  []*Theme  `json:"themes,omitempty"`
	Reviews []*Review `json:"reviews,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface for the Package type.
func (p *Package) MarshalJSON() ([]byte, error) {
	p.Version = CurrentVersion
	return json.Marshal(*p)
}
