package main

import (
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

	if (*add && *remove) || (!*add && !*remove) {
		log.Fatalln("improper usage")
	}
}
