package fb

import (
	"encoding/json"
	"sort"
	"sync/atomic"

	"github.com/pkg/errors"
)

// Attachment represents a Couch/PouchDB attachment.
type Attachment struct {
	refcount    int32
	ContentType string `json:"content_type"`
	Content     []byte `json:"data"`
}

// FileCollection represents a collection of Attachments which may be used by
// multiple related sub-document elements.
type FileCollection struct {
	files map[string]*Attachment
	views []*FileCollectionView
}

// GetFile returns an Attachment based on the file name. If the file does not
// the second return value will be false.
func (fc *FileCollection) GetFile(name string) (*Attachment, bool) {
	att, ok := fc.files[name]
	return att, ok
}

// FileCollectionView represents a view into a larger FileCollection, which can
// be used by sub-elements.
type FileCollectionView struct {
	col     *FileCollection
	members map[string]*Attachment
}

// NewFileCollection returns a new, empty FileCollection.
func NewFileCollection() *FileCollection {
	return &FileCollection{
		files: make(map[string]*Attachment),
		views: make([]*FileCollectionView, 0, 1),
	}
}

// AddView creates a new View on a FileCollection, which can be used by sub-elements.
func (fc *FileCollection) AddView(v *FileCollectionView) error {
	for filename := range v.members {
		att, ok := fc.files[filename]
		if !ok {
			return errors.New(filename + " not found in collection")
		}
		v.members[filename] = att
		atomic.AddInt32(&att.refcount, 1)
	}
	v.col = fc
	fc.views = append(fc.views, v)
	return nil
}

// RemoveView removes a FileCollectionView from a FileCollection
func (fc *FileCollection) RemoveView(v *FileCollectionView) error {
	for filename := range v.members {
		att, _ := fc.files[filename]
		atomic.AddInt32(&att.refcount, 1)
		if att.refcount == 0 {
			delete(fc.files, filename)
		}
	}
	for i, view := range fc.views {
		if view == v {
			fc.views = append(fc.views[:i], fc.views[i+1:]...)
			return nil
		}
	}
	return errors.New("Didn't find the view")
}

// NewView returns a new FileCollectionView from the existing FileCollection.
func (fc *FileCollection) NewView() *FileCollectionView {
	v := &FileCollectionView{
		col:     fc,
		members: make(map[string]*Attachment),
	}
	fc.views = append(fc.views, v)
	return v
}

// RemoveAll removes all references to the named Attachment.
func (fc *FileCollection) RemoveAll(name string) {
	delete(fc.files, name)
	for _, view := range fc.views {
		delete(view.members, name)
	}
}

// MarshalJSON implements the json.Marshaler interface for the FileCollection type.
func (fc *FileCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(fc.files)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the FileCollection type.
func (fc *FileCollection) UnmarshalJSON(data []byte) error {
	fc.files = make(map[string]*Attachment)
	fc.views = make([]*FileCollectionView, 0)
	if err := json.Unmarshal(data, &fc.files); err != nil {
		return errors.Wrap(err, "Unmarshal FileCollection")
	}
	return nil
}

// SetFile sets the requested attachment, replacing it if it already exists.
func (v *FileCollectionView) SetFile(name, ctype string, content []byte) {
	att := &Attachment{
		refcount:    1,
		ContentType: ctype,
		Content:     content,
	}
	v.col.files[name] = att
	v.members[name] = att
}

// AddFile adds the requested attachment. Returns an error if it already exists.
func (v *FileCollectionView) AddFile(name, ctype string, content []byte) error {
	if _, ok := v.col.files[name]; ok {
		return errors.New("File of that name already exists in the collection")
	}
	v.SetFile(name, ctype, content)
	return nil
}

// RemoveFile removes the named attachment from the collection.
func (v *FileCollectionView) RemoveFile(name string) error {
	att, ok := v.members[name]
	if !ok {
		return errors.New("File does not exist in view")
	}
	delete(v.members, name)
	atomic.AddInt32(&att.refcount, -1)
	if att.refcount == 0 {
		v.col.RemoveAll(name)
	}
	return nil
}

// FileList returns a list of filenames contained within the view.
func (v *FileCollectionView) FileList() []string {
	files := make([]string, 0, len(v.members))
	for name := range v.members {
		files = append(files, name)
	}
	return files
}

// GetFile returns an Attachment based on the file name. If the file does not
// the second return value will be false.
func (v *FileCollectionView) GetFile(name string) (*Attachment, bool) {
	att, ok := v.members[name]
	return att, ok
}

// MarshalJSON implements the json.Marshaler interface for the FileCollectionView type.
func (v *FileCollectionView) MarshalJSON() ([]byte, error) {
	names := make([]string, 0, len(v.members))
	for name := range v.members {
		names = append(names, name)
	}
	sort.Strings(names) // For consistent output
	return json.Marshal(names)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the FileCollectionView type.
func (v *FileCollectionView) UnmarshalJSON(data []byte) error {
	v.members = make(map[string]*Attachment)
	var names []string
	if err := json.Unmarshal(data, &names); err != nil {
		return errors.Wrap(err, "Unmarshal FileCollectionView")
	}
	for _, filename := range names {
		v.members[filename] = nil
	}
	return nil
}
