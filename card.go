package fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CardQueue int

const (
	QueueNew CardQueue = iota
	QueueLearning
	QueueReview
)

type Card struct {
	id       string
	Rev      *string
	Created  time.Time
	Modified time.Time
	Imported *time.Time
	// 	Queue       CardQueue
	// 	Suspended   bool
	// 	Buried      bool
	// 	AutoBuried  bool
	// 	Due         *time.Time
	// 	Interval    *time.Duration
	// 	SRSFactor   float32
	// 	ReviewCount int
	// 	LapseCount  int
}

type cardDoc struct {
	Type     string     `json:"type"`
	ID       string     `json:"_id"`
	Rev      *string    `json:"_rev,omitempty"`
	Created  time.Time  `json:"created"`
	Modified time.Time  `json:"modified"`
	Imported *time.Time `json:"imported,omitempty"`
	// 	Queue       CardQueue      `json:"state"`
	// 	Suspended   *bool          `json:"suspended,omitempty"`
	// 	Buried      *bool          `json:"buried,omitempty"`
	// 	AutoBuried  *bool          `json:"autoBuried,omitempty"`
	// 	Due         *time.Time     `json:"due,omitempty"`
	// 	Interval    *time.Duration `json:"interval,omitempty"`
	// 	SRSFactor   *float32       `json:"srsFactor,omitempty"`
	// 	ReviewCount *int           `json:"reviewCount,omitempty"`
	// 	LapseCount  *int           `json:"lapseCount,omitempty"`
}

func NewCard(noteID string, template int) (*Card, error) {
	return &Card{
		id: fmt.Sprintf("%s.%d", noteID, template),
	}, nil
}

func (c *Card) MarshalJSON() ([]byte, error) {
	doc := cardDoc{
		Type:     "card",
		ID:       "card-" + c.id,
		Rev:      c.Rev,
		Created:  c.Created,
		Modified: c.Modified,
		Imported: c.Imported,
		// 		Queue:    c.Queue,
		// 		Due:      c.Due,
		// 		Interval: c.Interval,
	}
	return json.Marshal(doc)
}

func (c *Card) UnmarshalJSON(data []byte) error {
	doc := &cardDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return err
	}
	if doc.Type != "card" {
		return errors.New("Invalid document type for card: " + doc.Type)
	}
	c.id = strings.TrimPrefix(doc.ID, "card-")
	c.Rev = doc.Rev
	c.Created = doc.Created
	c.Modified = doc.Modified
	c.Imported = doc.Imported
	// 	c.Queue = doc.Queue
	// 	if doc.Suspended != nil {
	// 		c.Suspended = *doc.Suspended
	// 	}
	// 	if doc.Buried != nil {
	// 		c.Buried = *doc.Buried
	// 	}
	// 	if doc.AutoBuried != nil {
	// 		c.AutoBuried = *doc.AutoBuried
	// 	}
	// 	c.Due = doc.Due
	// 	c.Interval = doc.Interval
	// 	if doc.SRSFactor != nil {
	// 		c.SRSFactor = *doc.SRSFactor
	// 	}
	// 	if doc.ReviewCount != nil {
	// 		c.ReviewCount = *doc.ReviewCount
	// 	}
	// 	if doc.LapseCount != nil {
	// 		c.LapseCount = *doc.LapseCount
	// 	}
	return nil
}

func (c *Card) Identity() string {
	return c.id
}

/*
func (c *Card) SetRev(rev string)        { c.Rev = &rev }
func (c *Card) DocID() string            { return "card-" + c.id }
func (c *Card) ImportedTime() *time.Time { return c.Imported }
func (c *Card) ModifiedTime() *time.Time { return &c.Modified }

func (c *Card) MergeImport(i interface{}) (bool, error) {
fmt.Printf("0\n")
	existing := i.(*Card)
	if c.id != existing.id {
		return false, errors.New("IDs don't match")
	}
	if !c.Created.Equal(existing.Created) {
		return false, errors.New("Created timestamps don't match")
	}
	c.Rev = existing.Rev
	if c.Modified.After(existing.Modified) {
		// The new version is newer than the existing one, so update
		return true, nil
	}
	// The new version is older, so we need to use the version we just read
	c.Modified = existing.Modified
	c.Imported = existing.Imported
	return false, nil
}
*/
