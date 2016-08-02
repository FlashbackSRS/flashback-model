package ankiconv

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/flimzy/anki"

	fb "github.com/flimzy/flashback-model"
)

func int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

type Bundle struct {
	apkg     *anki.Apkg
	b        *fb.Bundle
	now      *time.Time
	docs     []interface{}
	owner    *fb.User
	modelMap map[anki.ID]*fb.Model
}

func (bx *Bundle) MarshalJSON() ([]byte, error) {
	return json.Marshal(bx.docs)
}

func NewBundle() *Bundle {
	b := &Bundle{}
	now := time.Now()
	b.now = &now
	b.modelMap = make(map[anki.ID]*fb.Model)
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
	id := fb.KeyToIDString(bx.owner.UUID(), int64ToBytes(created.UnixNano()))
	b, _ := fb.NewBundle(id, o)
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
	if err := bx.addDecks(); err != nil {
		return fmt.Errorf("Error converting decks: %s", err)
	}
	if err := bx.addNotes(); err != nil {
		return fmt.Errorf("Error converting notes: %s", err)
	}
	// 	if err := bx.addCards(); err != nil {
	// 		return fmt.Errorf("Error converting cards: %s", err)
	// 	}
	return nil
}

type modelArray []*anki.Model

func (a modelArray) Len() int           { return len(a) }
func (a modelArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a modelArray) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (bx *Bundle) addThemes() error {
	c, err := bx.apkg.Collection()
	if err != nil {
		return err
	}
	models := modelArray(make([]*anki.Model, 0, len(c.Models)))
	for _, m := range c.Models {
		models = append(models, m)
	}
	sort.Sort(models)
	for _, aModel := range models {
		t, err := bx.convertTheme(aModel)
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
	id := fb.KeyToIDString(bx.owner.UUID(), int64ToBytes(int64(aModel.ID)))
	t, _ := fb.NewTheme(id)
	t.Name = &aModel.Name
	t.Modified = &modified
	t.Imported = bx.now
	t.SetFile("$main.css", "text/css", []byte(aModel.CSS))
	m, _ := t.NewModel(fb.ModelType(aModel.Type))
	m.Name = &aModel.Name
	bx.modelMap[aModel.ID] = m
	tNames := make([]string, len(aModel.Templates))
	for i, tmpl := range aModel.Templates {
		// TODO: store template names/order
		qName := "!" + aModel.Name + "." + tmpl.Name + " question.html"
		aName := "!" + aModel.Name + "." + tmpl.Name + " answer.html"
		m.AddFile(qName, fb.HTMLTemplateContentType, []byte(tmpl.QuestionFormat))
		m.AddFile(aName, fb.HTMLTemplateContentType, []byte(tmpl.AnswerFormat))
		tNames[i] = tmpl.Name
	}
	fields := make(map[int]*anki.Field)
	for _, field := range aModel.Fields {
		fields[field.Ordinal] = field
	}
	for i := 0; i < len(fields); i++ {
		field, ok := fields[i]
		if !ok {
			return nil, errors.New("Anki field missing")
		}
		m.AddField(fb.AnkiField, field.Name)
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

type deckArray []*anki.Deck

func (a deckArray) Len() int           { return len(a) }
func (a deckArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a deckArray) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (bx *Bundle) addDecks() error {
	c, err := bx.apkg.Collection()
	if err != nil {
		return err
	}
	decks := deckArray(make([]*anki.Deck, 0, len(c.Decks)))
	for _, d := range c.Decks {
		decks = append(decks, d)
	}
	sort.Sort(decks)
	for _, aDeck := range decks {
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
	id := fb.KeyToIDString(bx.owner.UUID(), int64ToBytes(int64(aDeck.ID)))
	d, _ := fb.NewDeck(id)
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

func (bx *Bundle) addNotes() error {
	notes, err := bx.apkg.Notes()
	if err != nil {
		return err
	}
	for notes.Next() {
		if aNote, err := notes.Note(); err != nil {
			return err
		} else {
			n, err := bx.convertNote(aNote)
			if err != nil {
				return err
			}
			for i, value := range aNote.FieldValues {
				fv := n.GetFieldValue(i)
				fv.SetText(value)
				if files, err := bx.extractFiles(value); err != nil {
					return err
				} else {
					for _, file := range files {
						fv.AddFile(file.Filename, file.ContentType, file.Content)
					}
				}
			}
			bx.docs = append(bx.docs, n)
		}
	}
	return nil
}

type fileType struct {
	RE      *regexp.Regexp
	TypeMap map[string]string
}

var FileTypes []fileType = []fileType{
	fileType{
		RE: regexp.MustCompile(`src="(.*)"`),
		TypeMap: map[string]string{
			".jpg": "image/jpeg",
			".png": "image/png",
			".gif": "image/gif",
		},
	},
	fileType{
		RE: regexp.MustCompile(`\[sound:(.*)\]`),
		TypeMap: map[string]string{
			".mp3": "audo/mpeg",
			".ogg": "audio/ogg",
			".oga": "audio/ogg",
			".spx": "audio/ogg",
			".wav": "audio/wav",
			".3gp": "audio/3gpp",
		},
	},
}

type Att struct {
	Filename    string
	ContentType string
	Content     []byte
}

func (bx *Bundle) extractFiles(value string) ([]*Att, error) {
	files := make([]*Att, 0, 1)
	for _, ft := range FileTypes {
		extracted := ft.RE.FindAllStringSubmatch(value, -1)
		for _, matches := range extracted {
			for _, filename := range matches[1:] {
				ext := strings.ToLower(filepath.Ext(filename))
				cType, ok := ft.TypeMap[ext]
				if !ok {
					return []*Att{}, fmt.Errorf("Unable to determine content type for `%s`", filename)
				}
				content, err := bx.apkg.ReadMediaFile(filename)
				if err != nil {
					return []*Att{}, err
				}
				files = append(files, &Att{
					Filename:    filename,
					ContentType: cType,
					Content:     content,
				})
			}
		}
	}
	return files, nil
}

func (bx *Bundle) convertNote(aNote *anki.Note) (*fb.Note, error) {
	id := fb.KeyToIDString(bx.owner.UUID(), int64ToBytes(int64(aNote.ID)))
	n, _ := fb.NewNote(id, bx.modelMap[aNote.ModelID])
	n.ID, _ = fb.NewID("note", id)
	return n, nil
}

/*
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
/*
func (bx *Bundle) convertCard(aCard *anki.Card) (*fb.Card, error) {
// 	c := fb.NewCard(
	return nil, nil
}

*/
