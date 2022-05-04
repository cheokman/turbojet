package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
)

const (
	storagePath = "/.turbojet"
	storageFile = "storage.json"
)

type Storage struct {
	Providers []Provider `json:"providers"`
}

func NewStorage() Storage {
	return Storage{
		Providers: []Provider{},
	}
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

func (s *Storage) NewProvider(pn string) Provider {
	p, ok := s.GetProvider(pn)

	if !ok {
		p = NewProvider()
		p.Name = pn
		s.PutProvider(p)
	}
	return p
}

func (s *Storage) GetProvider(pn string) (Provider, bool) {
	for i, p := range s.Providers {
		if p.Name == pn {
			return s.Providers[i], true
		}
	}
	return Provider{Name: pn}, false
}

func (s *Storage) GetProviders() []Provider {
	return s.Providers
}

func (s *Storage) PutProvider(provider Provider) {
	for i, p := range s.Providers {
		if p.Name == provider.Name {
			s.Providers[i] = provider
			return
		}
	}
	s.Providers = append(s.Providers, provider)
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
