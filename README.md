# bookmark

store aliased links to improve taking notes

### suggested workflow

use github gists/projects/issues to take notes, use this tool for search with the
ability to open a bookmark from the terminal.

### usage

I like to add aliases to each `*.json` I use for a specific category.

```bash
alias guides="bookmark --path ~/bookmarks/guides.json"
alias meetings="meetings --path ~/bookmarks/meetings.json"
# ... etc
```

then, I can easily use each aliased command to add/modify/search for the bookmark.

```bash
guides --add "my bookmark" www.google.com
guides "my bookmark" # opens www.google.com
guides --delete "my bookmark"
```

an easy way to incorporate fuzzy search is via another alias, such as

```bash 
alias search="bookmark --path ~/bookmarks/ $(bookmark --path ~/bookmarks/ --all | fzf | cut -f 1)"
```

which will fuzzy search and open that bookmark


### Planned future features

- [ ] link bookmarks together (graph-like)
- [ ] tag a bookmark, filter by tag
- [ ] edit support for saved bookmarks (replaces add -> delete -> add)
- [ ] add global config file (yaml)
- [ ] add option to use Sqlite instead of json (option controlled by config file)
