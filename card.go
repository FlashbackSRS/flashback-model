package fb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/*
type CardQueue int

const (
	QueueNew CardQueue = iota
	QueueLearning
	QueueReview
)
*/

// Card represents a struct of card-related statistics and configuration.
type Card struct {
	ID         string     `json:"_id"`
	Rev        *string    `json:"_rev,omitempty"`
	Created    time.Time  `json:"created"`
	Modified   time.Time  `json:"modified"`
	Imported   *time.Time `json:"imported,omitempty"`
	LastReview *time.Time `json:"lastReview,omitempty"`
	ModelID    string     `json:"model"`
	// 	Queue       CardQueue      `json:"state"`
	Suspended bool `json:"suspended,omitempty"`
	// 	Buried      *bool          `json:"buried,omitempty"`
	// 	AutoBuried  *bool          `json:"autoBuried,omitempty"`
	Due         *Due      `json:"due,omitempty"`
	BuriedUntil *Due      `json:"buriedUntil,omitempty"`
	Interval    *Interval `json:"interval,omitempty"`
	EaseFactor  float32   `json:"easeFactor,omitempty"`
	ReviewCount int       `json:"reviewCount,omitempty"`
	// 	LapseCount  *int           `json:"lapseCount,omitempty"`
	Context interface{} `json:"context,omitempty"`
}

// Validate validates that all of the data in the card appears valid and self
// consistent. A nil return value means no errors were detected.
func (c *Card) Validate() error {
	if c.ID == "" {
		return errors.New("id required")
	}
	if _, _, _, err := c.parseID(); err != nil {
		return err
	}
	if c.Created.IsZero() {
		return errors.New("created time required")
	}
	if c.Modified.IsZero() {
		return errors.New("modified time required")
	}
	if c.ModelID == "" {
		return errors.New("model id required")
	}
	if !strings.HasPrefix(c.ModelID, "theme-") {
		return errors.New("invalid type in model ID")
	}
	return nil
}

func (c *Card) parseID() (bundleID string, noteID string, templateID uint32, err error) {
	if !strings.HasPrefix(c.ID, "card-") {
		return "", "", 0, errors.New("invalid ID type")
	}
	parts := strings.Split(strings.TrimPrefix(c.ID, "card-"), ".")
	if len(parts) != 3 {
		return "", "", 0, errors.New("invalid ID format")
	}
	template, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", 0, errors.Wrap(err, "invalid TemplateID")
	}
	return parts[0], parts[1], uint32(template), nil
}

// NewCard returns a new Card instance, with the requested id
func NewCard(theme string, model uint32, id string) (*Card, error) {
	nowTime := now()
	c := &Card{
		Created:  nowTime,
		Modified: nowTime,
		ID:       id,
		ModelID:  fmt.Sprintf("%s/%d", theme, model),
	}
	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation failure")
	}
	return c, nil
}

// To avoid loops when (un)marshaling
type cardAlias Card

type jsonCard struct {
	cardAlias
	Type      string `json:"type"`
	Suspended *bool  `json:"suspended,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface for the Card type.
func (c *Card) MarshalJSON() ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, errors.Wrap(err, "validation error")
	}
	doc := jsonCard{
		Type:      "card",
		cardAlias: cardAlias(*c),
	}
	if c.Suspended {
		doc.Suspended = &c.Suspended
	}
	return json.Marshal(doc)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Card type.
func (c *Card) UnmarshalJSON(data []byte) error {
	doc := &jsonCard{}
	if err := json.Unmarshal(data, doc); err != nil {
		return err
	}
	if doc.Type != "card" {
		return errors.New("invalid document type for card: " + doc.Type)
	}
	*c = Card(doc.cardAlias)
	if doc.Suspended != nil {
		c.Suspended = *doc.Suspended
	}
	return errors.Wrap(c.Validate(), "validation error")
}

// Identity returns the identity of the card as a string.
func (c *Card) Identity() string {
	return strings.TrimPrefix(c.ID, "card-")
}

// SetRev sets the Card's _rev attribute
func (c *Card) SetRev(rev string) { c.Rev = &rev }

// DocID returns the Card's _id attribute
func (c *Card) DocID() string { return c.ID }

// ImportedTime returns the Card's imported time, or nil
func (c *Card) ImportedTime() *time.Time { return c.Imported }

// ModifiedTime returns the Card's last modified time
func (c *Card) ModifiedTime() *time.Time { return &c.Modified }

// MergeImport attempts to merge i into c, returning true on success, or false
// if no merge was necessary.
func (c *Card) MergeImport(i interface{}) (bool, error) {
	existing, ok := i.(*Card)
	if !ok {
		return false, errors.Errorf("i is %T, not *fb.Card", i)
	}
	if c.Identity() != existing.Identity() {
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

// BundleID returns the card's BundleID
func (c *Card) BundleID() string {
	bundleID, _, _, _ := c.parseID()
	return "bundle-" + bundleID
}

// TemplateID returns the card's TemplateID
func (c *Card) TemplateID() uint32 {
	_, _, templateID, _ := c.parseID()
	return templateID
}

// NoteID returns the card's NoteID
func (c *Card) NoteID() string {
	_, noteID, _, _ := c.parseID()
	return "note-" + noteID
}
