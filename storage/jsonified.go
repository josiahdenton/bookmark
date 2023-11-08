package storage

import (
	"encoding/json"
	"errors"
	"github.com/josiahdenton/bookmark/bookmark"
	"log"
	"os"
)

type JsonStorage struct {
	path    string
	aliases map[string]string
	ready   bool
}

func NewJson(path string) JsonStorage {
	return JsonStorage{path, make(map[string]string), false}
}

func (store *JsonStorage) Connect() error {
	content, err := os.ReadFile(store.path)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("creating w%v", err)
		fp, err := os.Create(store.path)
		if err != nil {
			log.Fatalf("failed to create a new file: %v", err)
		}
	} else if err != nil {
		log.Fatalf("failed to open a file: %v", err)
	}

	var aliases map[string]string
	err = json.Unmarshal(content, &aliases)
	if err != nil {
		log.Fatalf("failed to parse storage file: %v", err)
	}
	store.ready = true

	return nil
}

func (store *JsonStorage) Save(bookmark bookmark.Bookmark) error {
	if !store.ready {
		return ConnectionErr
	}

	store.aliases[bookmark.Alias] = bookmark.Url
	err := store.write()
	if err != nil {
		log.Fatalf("failed to save changes: %v", err)
	}

	return nil
}

func (store *JsonStorage) write() error {
	bytes, err := json.Marshal(store.aliases)
	if err != nil {
		log.Fatalf("faild to marshal alias map: %v", err)
	}
	err = os.WriteFile(store.path, bytes, 0666)
	if err != nil {
		log.Fatalf("failed so save alias map to file: %v", err)
	}

	return nil
}

func (store *JsonStorage) Find(alias string) (bookmark.Bookmark, error) {
	if url, ok := store.aliases[alias]; !ok {
		return bookmark.Bookmark{}, errors.New("no bookmark found")
	} else {
		return bookmark.Bookmark{Alias: alias, Url: url}, nil
	}
}

func (store *JsonStorage) Delete(alias string) error {
	if _, exists := store.aliases[alias]; !exists {
		return errors.New("alias does not exist")
	}
	delete(store.aliases, alias)
	return nil
}
