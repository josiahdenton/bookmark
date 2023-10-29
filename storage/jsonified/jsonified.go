package jsonified

import "github.com/josiahdenton/bookmark/bookmark"

type JsonStorage struct{}

func New() *JsonStorage {
	return &JsonStorage{}
}

func (js *JsonStorage) Save(bookmark bookmark.Bookmark) error {
	return nil
}

func (js *JsonStorage) Find(alias string) (bookmark.Bookmark, error) {
	return bookmark.Bookmark{}, nil
}

func (js *JsonStorage) Delete(alias string) (bool, error) {
	return false, nil
}
