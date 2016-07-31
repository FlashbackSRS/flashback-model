package fbmodel

import (
	"encoding/json"
	"time"
)

type Deck struct {
	ID
	Rev         *string
	Created     *time.Time
	Modified    *time.Time
	Imported    *time.Time
	Name        *string
	Description *string
}

type deckDoc struct {
	Type        string     `json:"type"`
	ID          ID         `json:"_id"`
	Rev         *string    `json:"_rev,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
	Modified    *time.Time `json:"modified,omitempty"`
	Imported    *time.Time `json:"imported,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
}

func NewDeck(id string) (*Deck, error) {
	d := &Deck{}
	did, err := NewID("deck", id)
	if err != nil {
		return nil, err
	}
	d.ID = did
	return d, nil
}

func (d *Deck) MarshalJSON() ([]byte, error) {
	return json.Marshal(deckDoc{
		Type:        "deck",
		ID:          d.ID,
		Rev:         d.Rev,
		Created:     d.Created,
		Modified:    d.Modified,
		Imported:    d.Imported,
		Name:        d.Name,
		Description: d.Description,
	})
}
