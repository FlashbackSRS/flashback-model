package ankiconv

import (
	"bytes"
	// 	"crypto/sha1"
	"encoding/binary"
	// 	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/flimzy/anki"

	fb "github.com/flimzy/flashback-model"
)

func int64ToByte(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

type Bundle struct {
	apkg  *anki.Apkg
	b     *fb.Bundle
	now   *time.Time
	docs  []interface{}
	owner *fb.User
}

func (bx *Bundle) id(in []byte) []byte {
	b := make([]byte, 0, len(in)+24)
	b = append(b, bx.owner.UUID()...)
	b = append(b, []byte("-anki-")...)
	b = append(b, in...)
	return b
	/*
		return generateID(bx.owner.UUID(), []byte("-anki-"), in)*/
}

func (bx *Bundle) MarshalJSON() ([]byte, error) {
	return json.Marshal(bx.docs)
}

func NewBundle() *Bundle {
	b := &Bundle{}
	now := time.Now()
	b.now = &now
	return b
}

func (bx *Bundle) SetNow(now time.Time) {
	bx.now = &now
}

func Convert(name string, o *fb.User, a *anki.Apkg) (*Bundle, error) {
	b := NewBundle()
	return b, b.Convert(name, o, a)
}

func (bx *Bundle) Convert(name string, o *fb.User, a *anki.Apkg) error {
	bx.apkg = a
	bx.owner = o
	c, err := bx.apkg.Collection()
	if err != nil {
		return err
	}
	created := time.Time(*c.Created)
	modified := time.Time(*c.Modified)
	b := fb.CreateBundle(bx.id(int64ToByte(created.UnixNano())), o)
	b.Created = &created
	b.Modified = &modified
	b.Name = &name
	b.Imported = bx.now
	docs := make([]interface{}, 0, 100)
	docs = append(docs, b)
	bx.docs = docs
	bx.b = b
	if err := bx.addThemes(); err != nil {
		return fmt.Errorf("Error converting themes: %s", err)
	}
	// 	if err := bx.addDecks(); err != nil {
	// 		return fmt.Errorf("Error converting decks: %s", err)
	// 	}
	// 	if err := bx.addCards(); err != nil {
	// 		return fmt.Errorf("Error converting cards: %s", err)
	// 	}
	return nil
}

func (bx *Bundle) addThemes() error {
	c, err := bx.apkg.Collection()
	if err != nil {
		return err
	}
	for _, model := range c.Models {
		t, err := bx.convertTheme(model)
		if err != nil {
			return err
		}
		bx.docs = append(bx.docs, t)
	}
	return nil
}

// TODO: convert the following fields
// 	Tags           []string          `json:"tags"`  // Anki saves the tags of the last added note to the current model
// 	Fields         []*Field          `json:"flds"`  // Array of Field objects
// 	SortField      int               `json:"sortf"` // Integer specifying which field is used for sorting in the browser
// 	Type           ModelType         `json:"type"`      // Model type: Standard or Cloze

func (bx *Bundle) convertTheme(aModel *anki.Model) (*fb.Theme, error) {
	modified := time.Time(*aModel.Modified)
	id := bx.id(int64ToByte(int64(aModel.ID)))
	t := fb.CreateTheme(id)
	t.Name = &aModel.Name
	t.Modified = &modified
	t.Imported = bx.now
	t.SetFile("$main.css", "text/css", []byte(aModel.CSS))
	m, err := t.NewModel(t.ID.Identity())
	if err != nil {
		return nil, err
	}
	m.Modified = &modified
	m.Imported = bx.now
	tNames := make([]string, len(aModel.Templates))
	for i, tmpl := range aModel.Templates {
		qName := "!" + aModel.Name + "." + tmpl.Name + " question.html"
		aName := "!" + aModel.Name + "." + tmpl.Name + " answer.html"
		m.AddFile(qName, fb.HTMLTemplateContentType, []byte(tmpl.QuestionFormat))
		m.AddFile(aName, fb.HTMLTemplateContentType, []byte(tmpl.AnswerFormat))
		tNames[i] = tmpl.Name
	}

	buf := new(bytes.Buffer)
	if err := masterTmpl.Execute(buf, tNames); err != nil {
		return nil, err
	}
	m.AddFile("$template.0.html", fb.HTMLTemplateContentType, buf.Bytes())
	return t, nil
}

var masterTmpl = template.Must(template.New("template.html").Delims("[[", "]]").Parse(`
{{ $g := . }}
[[- range $i, $Name := . ]]
	<div class="question" data-id="[[ $i ]]">
		{{template "[[ $Name ]] question.html" $g}}
	</div>
	<div class="answer" data-id="[[ $i ]]">
		{{template "[[ $Name ]] answer.html" $g}}
	</div>
[[ end -]]
`))

/*
func (bx *Bundle) addDecks() error {
	c, err := bx.apkg.Collection()
	if err != nil {
		return err
	}
	for _, aDeck := range c.Decks {
		d, err := bx.convertDeck(aDeck)
		if err != nil {
			return err
		}
		bx.docs = append(bx.docs, d)
	}
	return nil
}

// TODO: Conert the following fields
// type Deck struct {
//     ExtendedNewCardLimit    int              `json:"extendedNew"`      // Extended new card limit for custom study
//     ExtendedReviewCardLimit int              `json:"extendedRev"`      // Extended review card limit for custom study
//     ConfigID                ID               `json:"conf"`             // ID of option group from dconf in `col` table
//     NewToday                [2]int           `json:"newToday"`         // two number array used somehow for custom study
//     ReviewsToday            [2]int           `json:"revToday"`         // two number array used somehow for custom study
//     LearnToday              [2]int           `json:"lrnToday"`         // two number array used somehow for custom study
//     TimeToday               [2]int           `json:"timeToday"`        // two number array used somehow for custom study (in ms)
//     Config                  *DeckConfig      `json:"-"`
// }

func (bx *Bundle) convertDeck(aDeck *anki.Deck) (*fb.Deck, error) {
	d := fb.NewDeck(bx.ankiID(aDeck.ID))
	modified := time.Time(*aDeck.Modified)
	d.Modified = &modified
	if aDeck.Name != "" {
		d.Name = &aDeck.Name
	}
	if aDeck.Description != "" {
		d.Description = &aDeck.Description
	}
	return d, nil
}

func (bx *Bundle) addCards() error {
	cards, err := bx.apkg.Cards()
	if err != nil {
		return err
	}
	for cards.Next() {
		if aCard, err := cards.Card(); err != nil {
			return err
		} else {
			c, err := bx.convertCard(aCard)
			if err != nil {
				return err
			}
			bx.docs = append(bx.docs, c)
		}
	}
	return nil
}

func (bx *Bundle) convertCard(aCard *anki.Card) (*fb.Card, error) {
// 	c := fb.NewCard(
	return nil, nil
}*/
