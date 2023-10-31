package main

import (
	"fmt"
	"github.com/josiahdenton/bookmark/bookmark"
	"github.com/josiahdenton/bookmark/storage"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

func main() {
	// bookmark --add <alias> <url>
	// bookmark --pah bookmark.json --add <alias> <url>
	// bookmark --rmv <alias>
	// bookmark --edit <alias> (future - would not be too hard...)

	add := flag.BoolP("add", "a", false, "will add the alias and url to the bookmark")
	remove := flag.BoolP("rmv", "r", false, "will remove the alias and url from bookmark")
	preferPath := flag.StringP("path", "p", "", "override the default storage location")
	// TODO add edit flag

	flag.Parse()

	arguments := flag.Args()

	fmt.Println(arguments)

	if (*add && *remove) || (!*add && !*remove) {
		log.Fatalln("improper usage")
	}

	var store storage.JsonStorage
	if len(*preferPath) == 0 {
		path, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get user home dir: %v", err)
		}
		path += "/bookmarks.json"
		store = storage.NewJson(path)
	} else {
		store = storage.NewJson(*preferPath)
	}

	err := store.Connect()
	if err != nil {
		log.Fatalf("failed to connect to storage: %v", err)
	}
	action := bookmark.NewAction(&store)
	fmt.Println(action)

	switch {
	case *add:
		if len(arguments) < 2 {
			log.Fatalln("not enough arguments")
		}
		action.Save(bookmark.Bookmark{
			Alias: arguments[0],
			Url:   arguments[1],
		})
		fmt.Println("Successfully saved bookmark")
		break
	case *remove:
		fmt.Println("remove")
		break
	}
}
