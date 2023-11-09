package bookmarks

import "fmt"

type Bookmark struct {
	Alias string `json:"alias"`
	Url   string `json:"url"`
	Added string `json:"added"`
}

type Bookmarks struct {
	Active []Bookmark `json:"active"`
}

func New() Bookmarks {
	return Bookmarks{
		Active: make([]Bookmark, 0, 1),
	}
}

func (b Bookmark) String() string {
    return fmt.Sprintf("alias: %s, url: %s", b.Alias, b.Url)
}
