package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/josiahdenton/bookmark/bookmarks"
	"golang.org/x/exp/slices"
)

type JsonStorage struct {
	path      string
	bookmarks bookmarks.Bookmarks
	ready     bool
	readOnly  bool
}

func New(path string) JsonStorage {
	return JsonStorage{path, bookmarks.Bookmarks{}, false, false}
}

func (store *JsonStorage) Connect() error {
	if !store.pathIsOfTypeJson() {
		return fmt.Errorf("type of file node specified is not json or dir")
	}

	fi, err := os.Stat(store.path)
	if os.IsNotExist(err) {
		err = store.setupEmptyStorageFile()
		if err != nil {
			return fmt.Errorf("failed to setup storage file")
		}
		// retry if no error
		store.Connect()
	} else if err != nil {
		return fmt.Errorf("failed to stat file")
	}

	if fi.IsDir() {
		return store.connectToDir()
	}
	return store.connectToFile()
}

func (store *JsonStorage) connectToFile() error {
	content, err := os.ReadFile(store.path)
	if err != nil {
		return fmt.Errorf("failed to open storage file: %w", err)
	}

	store.bookmarks, err = store.parse(content)
	store.ready = true

	return nil
}

func (store *JsonStorage) connectToDir() error {
	store.bookmarks = bookmarks.New()
	err := filepath.WalkDir(store.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk: %w", err)
		}
		split := strings.Split(d.Name(), ".")
		fileType := split[len(split)-1]
		// only parse json files
		if fileType == "json" {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to open storage file: %w", err)
			}
			bms, err := store.parse(content)
			if err != nil {
				return err
			}
			store.bookmarks.Active = append(store.bookmarks.Active, bms.Active...)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk dir: %w", err)
	}

	return nil
}

func (store *JsonStorage) parse(content []byte) (bookmarks.Bookmarks, error) {
	var bms bookmarks.Bookmarks
	err := json.Unmarshal(content, &bms)
	if err != nil {
		return bookmarks.New(), fmt.Errorf("failed to parse storage file: %w", err)
	}
	return bms, nil
}

func (store *JsonStorage) setupEmptyStorageFile() error {
	fp, err := os.Create(store.path)
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
	log.Println("successfully created new store file")
	return nil
}

func (store *JsonStorage) pathIsOfTypeJson() bool {
	splits := strings.Split(store.path, ".")
	if splits[len(splits)-1] == "json" {
		return true
	}
	return false
}

func (store *JsonStorage) Save(bookmark bookmarks.Bookmark) error {
	if !store.ready {
		return ConnectionErr
	}
	if store.readOnly {
		return fmt.Errorf("cannot modify with dir path")
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
	if store.readOnly {
		return fmt.Errorf("cannot modify with dir path")
	}
	if len(store.bookmarks.Active) == 0 {
		return errors.New("no active bookmarks to delete")
	}

	for current, bookmark := range store.bookmarks.Active {
		if bookmark.Alias == alias {
			store.bookmarks.Active = slices.Delete(store.bookmarks.Active, current, current+1)
		}
	}

	return store.write()
}

func (store *JsonStorage) All() []bookmarks.Bookmark {
	return store.bookmarks.Active
}
