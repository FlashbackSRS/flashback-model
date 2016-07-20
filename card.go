package fbmodel

import ()

type Card struct {
	BaseDoc
	DeckID string `json:"deckId"`
}

func NewCard(id string, deckID string) *Card {
	c := &Card{}
	c.ID = id
	c.DeckID = deckID
	c.Type = "card"
	return c
}

// type Card struct {
//     NoteID         ID               `db:"nid"`    // Foreign Key to a Note
//     TemplateID     int              `db:"ord"`    // The Template ID, within the Model, to which this card corresponds.
//     Modified       TimestampSeconds `db:"mod"`    // Last modified time
//     UpdateSequence int              `db:"usn"`    // Update sequence number
//     Type           CardType         `db:"type"`   // Card type: new, learning, due
//     Queue          CardQueue        `db:"queue"`  // Queue: suspended, user buried, sched buried
//     Due            TimestampSeconds `db:"due"`    // Time when the card is next due
//     Interval       DurationSeconds  `db:"ivl"`    // SRS interval in seconds
//     Factor         float32          `db:"factor"` // SRS factor
//     ReviewCount    int              `db:"reps"`   // Number of reviews
//     Lapses         int              `db:"lapses"` // Number of times card went from "answered correctly" to "answered incorrectly" state
//     Left           int              `db:"left"`   // Reviews remaining until graduation
//     OriginalDue    TimestampSeconds `db:"odue"`   // Original due time. Only used when card is in filtered deck.
//     OriginalDeckID ID               `db:"odid"`   // Original Deck ID. Only used when card is in filtered deck.
// }
