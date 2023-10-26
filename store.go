package main

import (
	"encoding/json"
	"errors"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	// store --add <alias> <url>
	// store --rmv <alias>
	// store --edit <alias> (future - would not be too hard...)

	add := flag.BoolP("add", "a", false, "will add the alias and url to the store")
	remove := flag.BoolP("rmv", "r", false, "will remove the alias and url from store")
	// TODO add edit flag

	if (*add && *remove) || (!*add && !*remove) {
		log.Fatalln("improper usage")
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to find home dir: %v", err)
	}
	storePath := path.Join(homePath, ".store.json")

	store := fill(storePath)
}

type Store struct {
	Aliases []Alias `json:"aliases"`
}

type Alias struct {
	Alias string `json:"alias"`
	Url   string `json:"url"`
}

func fill(path string) Store {
	var fp *os.File
	fp, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		fp, err = os.Create(path)
		if err != nil {
			log.Fatalf("failed to create initial store file: %v", err)
		}
	} else if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}
	// fp is set
	// read in to json
	var buffer [100]byte
	content := make([]byte, 0, 100)

	for {
		_, err = fp.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Fatalf("failed to read store file: %v", err)
		}
	}

	var store Store
	json.Unmarshal(&store)
}
