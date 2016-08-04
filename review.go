package fb

import (
	"time"
)

type Review struct {
	CardID           string         `json:"cardID"`
	Timestamp        *time.Time     `json:"timestamp"`
	Ease             ReviewEase     `json:"ease"`
	Interval         *time.Duration `json:"interval"`
	PreviousInterval *time.Duration `json:"previousInterval"`
	SRSFactor        float32        `json:"srsFactor"`
	ReviewTime       *time.Duration `json:"reviewTime"`
	Type             ReviewType     `json:"reviewType"`
}

type ReviewEase int

const (
	ReviewEaseWrong ReviewEase = 1
	ReviewEaseHard  ReviewEase = 2
	ReviewEaseOK    ReviewEase = 3
	ReviewEaseEasy  ReviewEase = 4
)

type ReviewType int

const (
	ReviewTypeLearn ReviewType = iota
	ReviewTypeReview
	ReviewTypeRelearn
	ReviewTypeCram
)

func NewReview(cardID string) (*Review, error) {
	return &Review{
		CardID: cardID,
	}, nil
}
