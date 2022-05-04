package content

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"turbojet/cli"
)

const (
	storagePath = "/.turbojet"
	storageFile = "content.json"
)

type Storage struct {
	PropertyContents []PropertyContent `json:"property_contents"`
	Contents         []Content         `json:"contents"`
}

func NewStorage() Storage {
	return Storage{
		PropertyContents: []PropertyContent{},
	}
}

func (s *Storage) GetContentByID(ID string) (Content, bool) {
	for i, c := range s.Contents {
		if c.ID == ID {
			return s.Contents[i], true
		}
	}
	return Content{ID: ID}, false
}

func (s *Storage) PutContent(content Content) {
	for i, c := range s.Contents {
		if content.From == fromLocal && c.From == fromLocal {
			if c.ID == content.ID {
				s.Contents[i] = content
				return
			}
		}
	}
	s.Contents = append(s.Contents, content)
}

func LoadStorage(path string, w io.Writer) (Storage, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewStorage(), nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return NewStorage(), fmt.Errorf("reading provider config from '%s' failed %v", path, err)
	}

	return NewStorageFromBytes(bytes)
}

func SaveStorage(s Storage) error {
	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	path := GetStoragePath() + string(os.PathSeparator) + storageFile
	err = ioutil.WriteFile(path, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func NewStorageFromBytes(bytes []byte) (Storage, error) {
	var s Storage
	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func GetStoragePath() string {
	path := GetHomePath() + storagePath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	return path
}

func GetHomePath() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

func LoadContentFromStorageByID(c *cli.Context, path string, id string) (Content, Storage, bool) {
	contexts, storage, err := LoadContentsFromStorage(c, path)
	if err != nil {
		return Content{}, storage, false
	}
	for i, content := range contexts {
		if content.ID == id {
			return contexts[i], storage, true
		}
	}
	return Content{}, storage, false
}

func LoadContentsFromStorage(c *cli.Context, path string) ([]Content, Storage, error) {
	w := c.Writer()
	storage, err := LoadStorage(path, w)
	if err != nil {
		return []Content{}, storage, err
	}

	return storage.Contents, storage, nil
}
