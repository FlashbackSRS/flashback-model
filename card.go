package fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CardState int

const (
	NewState CardState = iota
	LearnedState
	DueState
)

type Card struct {
	id          string
	Rev         *string
	Created     *time.Time
	Modified    *time.Time
	Imported    *time.Time
	State       CardState
	Suspended   bool
	Buried      bool
	AutoBuried  bool
	Due         *time.Time
	Interval    *time.Duration
	SRSFactor   float32
	ReviewCount int
	LapseCount  int
}

type cardDoc struct {
	Type        string         `json:"type"`
	ID          string         `json:"_id"`
	Rev         *string        `json:"_rev,omitempty"`
	Created     *time.Time     `json:"created,omitempty"`
	Modified    *time.Time     `json:"modified,omitempty"`
	Imported    *time.Time     `json:"imported,omitempty"`
	State       CardState      `json:"state"`
	Suspended   *bool          `json:"suspended,omitempty"`
	Buried      *bool          `json:"buried,omitempty"`
	AutoBuried  *bool          `json:"autoBuried,omitempty"`
	Due         *time.Time     `json:"due,omitempty"`
	Interval    *time.Duration `json:"interval,omitempty"`
	SRSFactor   *float32       `json:"srsFactor,omitempty"`
	ReviewCount *int           `json:"reviewCount,omitempty"`
	LapseCount  *int           `json:"lapseCount,omitempty"`
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
		State:    c.State,
		Due:      c.Due,
		Interval: c.Interval,
	}
	if c.Suspended {
		doc.Suspended = &c.Suspended
	}
	if c.Buried {
		doc.Buried = &c.Buried
	}
	if c.AutoBuried {
		doc.AutoBuried = &c.AutoBuried
	}
	if c.SRSFactor > 0 {
		doc.SRSFactor = &c.SRSFactor
	}
	if c.ReviewCount > 0 {
		doc.ReviewCount = &c.ReviewCount
	}
	if c.LapseCount > 0 {
		doc.LapseCount = &c.LapseCount
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
	c.State = doc.State
	if doc.Suspended != nil {
		c.Suspended = *doc.Suspended
	}
	if doc.Buried != nil {
		c.Buried = *doc.Buried
	}
	if doc.AutoBuried != nil {
		c.AutoBuried = *doc.AutoBuried
	}
	c.Due = doc.Due
	c.Interval = doc.Interval
	if doc.SRSFactor != nil {
		c.SRSFactor = *doc.SRSFactor
	}
	if doc.ReviewCount != nil {
		c.ReviewCount = *doc.ReviewCount
	}
	if doc.LapseCount != nil {
		c.LapseCount = *doc.LapseCount
	}
	return nil
}

func (c *Card) Identity() string {
	return c.id
}
