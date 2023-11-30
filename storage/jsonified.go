package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/josiahdenton/bookmark/bookmarks"
	"golang.org/x/exp/slices"
)

type JsonStorage struct {
	path      string
	bookmarks bookmarks.Bookmarks
	ready     bool
}

func New(path string) JsonStorage {
	return JsonStorage{path, bookmarks.Bookmarks{}, false}
}

func (store *JsonStorage) Connect() error {
	content, err := os.ReadFile(store.path)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("file does not exist: %v", err)
		err = setupEmptyStorageFile(store.path)
		if err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("failed to open storage file: %w", err)
	}

	var bookmarks bookmarks.Bookmarks
	err = json.Unmarshal(content, &bookmarks)
	if err != nil {
		return fmt.Errorf("failed to parse storage file: %w", err)
	}

	store.ready = true
	store.bookmarks = bookmarks
	return nil
}

func setupEmptyStorageFile(path string) error {
	log.Println("creating new storage file")
	fp, err := os.Create(path)
	defer fp.Close()
	if err != nil {
		return fmt.Errorf("failed to create a new file: %v", err)
	}

	empty := bookmarks.New()
	bytes, err := json.Marshal(empty)

	_, err = fp.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write empty store file: %v", err)
	}
	log.Println("successfully created new store file, please retry previous command")
	return nil
}

func (store *JsonStorage) Save(bookmark bookmarks.Bookmark) error {
	if !store.ready {
		return ConnectionErr
	}
	for _, val := range store.bookmarks.Active {
		if val == bookmark {
			return errors.New("duplicate bookmark id")
		}
	}
	store.bookmarks.Active = append(store.bookmarks.Active, bookmark)
	err := store.write()
	if err != nil {
		return fmt.Errorf("failed to save changes: %v", err)
	}

	return nil
}

func (store *JsonStorage) write() error {
	bytes, err := json.Marshal(store.bookmarks)
	if err != nil {
		return fmt.Errorf("failed to marshal bookmarks: %v", err)
	}
	err = os.WriteFile(store.path, bytes, 0666)
	if err != nil {
		return fmt.Errorf("failed so save alias map to file: %v", err)
	}

	return nil
}

func (store *JsonStorage) Find(alias string) (bookmarks.Bookmark, error) {
	for _, bookmark := range store.bookmarks.Active {
		if bookmark.Alias == alias {
			return bookmark, nil
		}
	}
	return bookmarks.Bookmark{}, errors.New("no bookmark found")
}

func (store *JsonStorage) Delete(alias string) error {
	if len(store.bookmarks.Active) == 0 {
		return errors.New("no active bookmarks to delete")
	}

	for current, bookmark := range store.bookmarks.Active {
		if bookmark.Alias == alias {
			slices.Delete(store.bookmarks.Active, current, current+1)
		}
	}

	return nil
}

func (store *JsonStorage) All() []bookmarks.Bookmark {
	return store.bookmarks.Active
}
