package fbmodel

import (
	"encoding/json"
	"errors"
)

type FileCollection struct {
	files map[string]*Attachment
	views []*FileCollectionView
}

type FileCollectionView struct {
	col     *FileCollection
	members map[string]*Attachment
}

func NewFileCollection() *FileCollection {
	return &FileCollection{
		files: make(map[string]*Attachment),
		views: make([]*FileCollectionView, 0, 1),
	}
}

func (fc *FileCollection) NewView() *FileCollectionView {
	v := &FileCollectionView{
		col:     fc,
		members: make(map[string]*Attachment),
	}
	fc.views = append(fc.views, v)
	return v
}

func (fc *FileCollection) removeFile(name string) {
	delete(fc.files, name)
	for _, view := range fc.views {
		delete(view.members, name)
	}
}

func (fc *FileCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(fc.files)
}

// Sets the requested attachment, replacing it if it already exists.
func (v *FileCollectionView) SetFile(name, ctype string, content []byte) {
	att := &Attachment{
		ContentType: ctype,
		Content:     content,
	}
	v.col.files[name] = att
	v.members[name] = att
}

// Adds the requested attachment. Returns an error if it already exists.
func (v *FileCollectionView) AddFile(name, ctype string, content []byte) error {
	if _, ok := v.col.files[name]; ok {
		return errors.New("File of that name already exists in the collection")
	}
	v.SetFile(name, ctype, content)
	return nil
}

func (v *FileCollectionView) RemoveFile(name string) {
	v.col.removeFile(name)
}

func (v *FileCollectionView) GetFile(name string) (*Attachment, bool) {
	att, ok := v.members[name]
	return att, ok
}

func (v *FileCollectionView) MarshalJSON() ([]byte, error) {
	names := make([]string, 0, len(v.members))
	for name, _ := range v.members {
		names = append(names, name)
	}
	return json.Marshal(names)
}
