package main

import (
	"fmt"
	"github.com/josiahdenton/bookmark/bookmark"
	"github.com/josiahdenton/bookmark/storage/jsonified"
	flag "github.com/spf13/pflag"
	"log"
)

func main() {
	// bookmark --add <alias> <url>
	// bookmark --rmv <alias>
	// bookmark --edit <alias> (future - would not be too hard...)

	add := flag.BoolP("add", "a", false, "will add the alias and url to the bookmark")
	remove := flag.BoolP("rmv", "r", false, "will remove the alias and url from bookmark")
	// TODO add edit flag

	flag.Parse()

	arguments := flag.Args()

	fmt.Println(arguments)

	if (*add && *remove) || (!*add && !*remove) {
		log.Fatalln("improper usage")
	}

	storage := jsonified.New()
	action := bookmark.NewAction(storage)
	fmt.Println(action)

	switch {
	case *add:
		//action.Save()
		fmt.Println("add")
		break
	case *remove:
		fmt.Println("remove")
		break
	}
}
