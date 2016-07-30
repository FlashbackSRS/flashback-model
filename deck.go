package fbmodel

import ()

type Deck struct {
	namedDoc
}

func NewDeck(id string) *Deck {
	d := &Deck{}
	d.doc = NewDoc("deck", id)
	return d
}

// func (d *Deck) NewNote(id string) *Note {
// 	n := &Note{}
// 	n.ID
// }
