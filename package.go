package fb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
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
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Bundle   *Bundle   `json:"bundle,omitempty"`
	Cards    []*Card   `json:"cards,omitempty"`
	Notes    []*Note   `json:"notes,omitempty"`
	Decks    []*Deck   `json:"decks,omitempty"`
	Themes   []*Theme  `json:"themes,omitempty"`
	Reviews  []*Review `json:"reviews,omitempty"`
}

type packageAlias Package

type jsonPackage struct {
	packageAlias
	Version int `json:"version"`
}

// MarshalJSON implements the json.Marshaler interface for the Package type.
func (p *Package) MarshalJSON() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	doc := jsonPackage{
		Version:      CurrentVersion,
		packageAlias: packageAlias(*p),
	}
	return json.Marshal(doc)
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (p *Package) UnmarshalJSON(data []byte) error {
	pa := &packageAlias{}
	if err := json.Unmarshal(data, &pa); err != nil {
		return err
	}
	*p = Package(*pa)
	modelMap := make(map[string]*Model)
	for _, t := range p.Themes {
		for _, m := range t.Models {
			modelMap[fmt.Sprintf("%s/%d", t.ID, m.ID)] = m
		}
	}
	for _, n := range p.Notes {
		key := fmt.Sprintf("%s/%d", n.ThemeID, n.ModelID)
		m, ok := modelMap[key]
		if !ok {
			return errors.Errorf("note %s has no matching model %s", n.ID, key)
		}
		n.Model = m
	}
	return nil
}

// Validate does some basic sanity checking on the package.
func (p *Package) Validate() error {
	cardMap := map[string]*Card{}
	for _, c := range p.Cards {
		if err := c.Validate(); err != nil {
			return errors.Wrapf(err, "card '%s' validation'", c.ID)
		}
		cardMap[c.ID] = c
	}

	cards := make([]*Card, 0, len(cardMap))

	for _, d := range p.Decks {
		if err := d.Validate(); err != nil {
			return errors.Wrapf(err, "deck '%s' validation", d.ID)
		}
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
