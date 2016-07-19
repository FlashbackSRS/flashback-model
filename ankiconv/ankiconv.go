package ankiconv

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/flimzy/anki"

	// 	"github.com/flimzy/flashback-model"
	"github.com/flimzy/flashback-model/bundle"
	"github.com/flimzy/flashback-model/theme"
	"github.com/flimzy/flashback-model/user"
)

func generateID(in ...[]byte) string {
	h := sha1.New()
	for _, part := range in {
		h.Write(part)
	}
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

type Bundle struct {
	b    *bundle.Bundle
	now  *time.Time
	docs []interface{}
}

func (b *Bundle) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.docs)
}

func Convert(name string, o *user.User, a *anki.Apkg) (*Bundle, error) {
	c, err := a.Collection()
	if err != nil {
		return nil, err
	}
	created := time.Time(*c.Created)
	modified := time.Time(*c.Modified)
	now := time.Now()
	id := generateID(o.UUID(), []byte("-anki-"), []byte(strconv.FormatInt(created.UnixNano(), 10)))
	fmt.Printf("Bundle ID = %s\n", id)
	b := bundle.New(id, o)
	b.Created = &created
	b.Modified = &modified
	b.Name = &name
	b.Imported = &now
	docs := make([]interface{}, 0, 100)
	docs = append(docs, b)
	bx := &Bundle{
		b:    b,
		now:  &now,
		docs: docs,
	}
	if err := bx.addThemes(c); err != nil {
		return nil, fmt.Errorf("Error processing themes: %s", err)
	}
	return bx, nil
}

func (bx *Bundle) addThemes(c *anki.Collection) error {
	for mid, m := range c.Models {
		t, err := bx.convertTheme(mid, m)
		if err != nil {
			return err
		}
		bx.docs = append(bx.docs, t)
	}
	return nil
}

func (bx *Bundle) convertTheme(mid anki.ID, model *anki.Model) (*theme.Theme, error) {
	modified := time.Time(*model.Modified)
	bmid := make([]byte, 8)
	binary.LittleEndian.PutUint64(bmid, uint64(mid))
	id := generateID(bx.b.Owner.UUID(), []byte("-anki-"), bmid)
	t := theme.NewTheme(id)
	t.Name = &model.Name
	t.Modified = &modified
	t.Imported = bx.now
	t.SetFile("$main.css", "text/css", []byte(model.CSS))
	m := t.NewModel(t.ID)
	m.Modified = &modified
	m.Imported = bx.now
	tNames := make([]string, len(model.Templates))
	for i, tmpl := range model.Templates {
		qName := "!" + model.Name + "." + tmpl.Name + " question.html"
		aName := "!" + model.Name + "." + tmpl.Name + " answer.html"
		m.AddFile(qName, theme.HTMLTemplateContentType, []byte(tmpl.QuestionFormat))
		m.AddFile(aName, theme.HTMLTemplateContentType, []byte(tmpl.AnswerFormat))
		tNames[i] = tmpl.Name
	}

	buf := new(bytes.Buffer)
	if err := masterTmpl.Execute(buf, tNames); err != nil {
		return nil, err
	}
	m.AddFile("$template.0.html", theme.HTMLTemplateContentType, buf.Bytes())

	return t, nil
}

// TODO: convert the following fields
// 	Tags           []string          `json:"tags"`  // Anki saves the tags of the last added note to the current model
// 	Fields         []*Field          `json:"flds"`  // Array of Field objects
// 	SortField      int               `json:"sortf"` // Integer specifying which field is used for sorting in the browser
// 	Type           ModelType         `json:"type"`      // Model type: Standard or Cloze

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
