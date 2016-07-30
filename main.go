package fbmodel

import (
	"encoding/json"
	"time"
)

const HTMLTemplateContentType = "text/html+flashbacktmpl"

type coreDoc struct {
	ID   string
	Type string
}

type DocID struct {
	coreDoc
}

func (d *DocID) String() string {
	return d.Type + "-" + d.ID
}

func (d *DocID) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

type DocType struct {
	coreDoc
}

func (d *DocType) String() string {
	return d.Type
}

func (d *DocType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func NewCoreDoc(doctype, id string) coreDoc {
	return coreDoc{
		ID: id,
		Type: doctype,
	}
}

type BaseDoc struct {
	doc      coreDoc
	ID       *DocID     `json:"_id"`
	Rev      *string    `json:"_rev,omitempty"`
	Type     *DocType   `json:"type"`
	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified"`
	Imported *time.Time `json:"imported,omitempty"`
}

func NewBaseDoc(doctype, id string) BaseDoc {
	cd := NewCoreDoc(doctype, id)
	doc := BaseDoc{}
	doc.doc = cd;
	doc.ID = &DocID{cd}
	doc.Type = &DocType{cd}
	return doc
}

type NamedDoc struct {
	BaseDoc
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func NewNamedDoc(doctype, id string) NamedDoc {
	nd := NamedDoc{}
	nd.BaseDoc = NewBaseDoc(doctype, id)
	return nd
}

type Attachment struct {
	ContentType string `json:"content-type"`
	Content     []byte `json:"data"`
}
