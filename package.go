package fb

import (
	"encoding/json"
	"fmt"
)

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

// Validate does some basic sanity checking on the package.
func (p *Package) Validate() error {
	cardMap := map[string]*Card{}
	for _, c := range p.Cards {
		cardMap[c.Identity()] = c
	}

	cards := make([]*Card, 0, len(cardMap))

	for _, d := range p.Decks {
		for _, id := range d.Cards.All() {
			c, ok := cardMap[id]
			if !ok {
				return fmt.Errorf("card '%s' listed in deck, but not found in package", id)
			}
			cards = append(cards, c)
			delete(cardMap, id)
		}
	}
	for id := range cardMap {
		return fmt.Errorf("card '%s' found in package, but not in a deck", id)
	}
	return nil
}
