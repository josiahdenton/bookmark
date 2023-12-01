package main

import (
	"fmt"
	"github.com/josiahdenton/bookmark/bookmarks"
	"github.com/josiahdenton/bookmark/storage"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

func main() {
	// bookmark --add <alias> <url>
	// bookmark --path bookmark.json --add <alias> <url>
	//     !!! for now, path would need to be included in every command...
	// bookmark --rmv <alias>
	// bookmark --edit <alias> (future - would not be too hard...)
	// bookmark --all

	// bookmark <alias>
	// above defaults to open case

	add := flag.BoolP("add", "a", false, "will add the alias and url to the bookmark")
	remove := flag.BoolP("rmv", "r", false, "will remove the alias and url from bookmark")
	all := flag.Bool("all", false, "dump all urls into stdout")
	preferPath := flag.StringP("path", "p", "", "override the default storage location")
	// TODO add edit flag

	flag.Parse()

	arguments := flag.Args()

	if (*add && *remove) || (*add && *all) || (*remove && *all) {
		log.Fatalln("improper usage")
	}

	var store storage.JsonStorage
	if len(*preferPath) == 0 {
		path, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get user home dir: %v", err)
		}
		path += "/bookmarks.json"
		store = storage.New(path)
	} else {
		store = storage.New(*preferPath)
	}

	err := store.Connect()
	if err != nil {
		log.Fatalf("failed to connect to storage: %v", err)
	}
	action := bookmarks.NewAction(&store)

	switch {
	case *add:
		mustHaveLength(arguments, 2)
		action.Save(bookmarks.Bookmark{
			Alias: arguments[0],
			Url:   arguments[1],
		}).Must()
		fmt.Println("Successfully saved bookmark")
		break
	case *remove:
		mustHaveLength(arguments, 1)
		action.Delete(arguments[0]).Must()
		break
	case *all:
		bookmarks := action.All()
		for _, bookmark := range bookmarks {
			fmt.Println(bookmark)
		}
		break
	default:
		mustHaveLength(arguments, 1)
		action.Find(arguments[0]).And().Open()
	}
}

func mustHaveLength(args []string, n int) {
	if len(args) < n {
		log.Fatalln("not enough arguments")
	}
}
