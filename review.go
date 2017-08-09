package fb

import (
	"errors"
	"strings"
	"time"
)

// Review represents a single card-review event.
type Review struct {
	CardID    string    `json:"cardID"`
	Timestamp time.Time `json:"timestamp"`
	// Ease             ReviewEase     `json:"ease"`
	// Interval         *time.Duration `json:"interval"`
	// PreviousInterval *time.Duration `json:"previousInterval"`
	// SRSFactor        float32        `json:"srsFactor"`
	// ReviewTime       *time.Duration `json:"reviewTime"`
	// Type             ReviewType     `json:"reviewType"`
}

// Validate validates that all of the data in the review appears valid and self
// consistent. A nil return value means no errors were detected.
func (r *Review) Validate() error {
	if r.CardID == "" {
		return errors.New("card id required")
	}
	if err := validateDocID(r.CardID); err != nil {
		return err
	}
	if !strings.HasPrefix(r.CardID, "card-") {
		return errors.New("incorrect doc type for card ID")
	}
	if r.Timestamp.IsZero() {
		return errors.New("timestamp required")
	}
	return nil
}

// type ReviewEase int
//
// const (
// 	ReviewEaseWrong ReviewEase = 1
// 	ReviewEaseHard  ReviewEase = 2
// 	ReviewEaseOK    ReviewEase = 3
// 	ReviewEaseEasy  ReviewEase = 4
// )
//
// type urReviewType int
//
// const (
// 	ReviewTypeLearn ReviewType = iota
// 	ReviewTypeReview
// 	ReviewTypeRelearn
// 	ReviewTypeCram
// )

// NewReview returns a new, empty Review for the provided Card.
func NewReview(cardID string) (*Review, error) {
	r := &Review{
		CardID:    cardID,
		Timestamp: now(),
	}
	return r, r.Validate()
}
