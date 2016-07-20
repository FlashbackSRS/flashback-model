package fbmodel

import (
	"time"
)

const HTMLTemplateContentType = "text/html+flashbacktmpl"

type BaseDoc struct {
	ID       string     `json:"_id"`
	Rev      string     `json:"_rev,omitempty"`
	Type     string     `json:"type"`
	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified"`
	Imported *time.Time `json:"imported,omitempty"`
}

type NamedDoc struct {
	BaseDoc
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type Attachment struct {
	ContentType string `json:"content-type"`
	Content     []byte `json:"data"`
}
