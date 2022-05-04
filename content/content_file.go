package content

import (
	"path/filepath"
)

type ContentFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Ext  string `json:"ext"`
	Size int    `json:"size"`
}

func NewContentFile(file string) (c *ContentFile, err error) {
	name, err := getName(file)
	if err != nil {
		return
	}
	c = &ContentFile{
		Name: name,
		Path: getPath(file),
		Ext:  getExt(file),
	}
	return
}

func getName(f string) (name string, err error) {
	name = filepath.Base(f)
	if name == "" {
		err = NewInvalidContentFileError(f)
	}
	return
}

func getPath(f string) string {
	return filepath.Dir(f)
}

func getExt(f string) string {
	return filepath.Ext(f)
}

func (c *ContentFile) SetSize(s int) {
	c.Size = s
}
