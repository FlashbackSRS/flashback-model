package fbmodel

import ()

type Deck struct {
	NamedDoc
}

func NewDeck(id string) *Deck {
	d := &Deck{}
	d.ID = id
	d.Type = "deck"
	return d
}
