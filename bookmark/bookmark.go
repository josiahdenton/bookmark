package bookmark

type Bookmark struct {
	Id    int    `json:"id"`
	Alias string `json:"alias"`
	Url   string `json:"url"`
	// tags??
}

type Bookmarks struct {
	Bookmarks []Bookmark `json:"bookmarks"`
}
