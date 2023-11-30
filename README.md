# bookmark

store aliased links to improve taking notes

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
alias search="bookmark --path ~/bookmarks/saved.json $(bookmark --path ~/bookmarks/saved.json --all | fzf | cut -f 1)"
```

which will fuzzy search and open that bookmark

### Bugs

- [ ] fix issue when file does not exist, make this a hidden error and not need an empty json file

### Features

- [ ] make dirs work for a path
