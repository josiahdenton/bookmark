package jsonified

import (
	"encoding/json"
	"errors"
	"github.com/josiahdenton/bookmark/bookmark"
	"log"
	"os"
)

type JsonStorage struct {
	path    string
	aliases Aliases
}

type Aliases struct {
	ActiveUrls []Alias `json:"active_urls"`
}

type Alias struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func New(path string) *JsonStorage {
	return &JsonStorage{path, Aliases{}}
}

func (storage *JsonStorage) Connect() error {
	//fp, err := os.Open(storage.path)
	content, err := os.ReadFile(storage.path)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("%v", err)
		content = retryRead(storage.path)
	} else if err != nil {
		log.Fatalf("failed to open a file: %v", err)
	}

	var aliases Aliases
	err = json.Unmarshal(content, &aliases)
	if err != nil {
		log.Fatalf("failed to parse storage file: %v", err)
	}

	storage.aliases = aliases

	return nil
}

func retryRead(path string) []byte {
	fp, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed to create a new file: %v", err)
	}
	err = fp.Close()
	if err != nil {
		log.Fatalf("failed to close created file: %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read: %v", err)
	}

	return content
}

func (storage *JsonStorage) Save(bookmark bookmark.Bookmark) error {
	return nil
}

func (storage *JsonStorage) Find(alias string) (bookmark.Bookmark, error) {
	return bookmark.Bookmark{}, nil
}

func (storage *JsonStorage) Delete(alias string) (bool, error) {
	return false, nil
}
