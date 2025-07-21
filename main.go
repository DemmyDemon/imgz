package main

import (
	"flag"
	"os"
	"path"
	"strings"

	"github.com/DemmyDemon/imgz/internal/do"

	"github.com/hajimehoshi/ebiten/v2"
)

var fileExts = []string{
	".jpg",
	".jpeg",
	".png",
	".webp",
	".gif",
}

func isRelevantFile(name string) bool {
	for _, ext := range fileExts {
		if strings.HasSuffix(strings.ToLower(name), ext) {
			return true
		}
	}
	return false
}
func init() {
	flag.BoolVar(&do.Verbosity, "verbose", false, "Enable verbose output")
	flag.Parse()

}

func main() {
	wd, err := os.Getwd()
	do.Fuck(err)

	if flag.NArg() >= 1 {
		wd = flag.Arg(0)
	}

	entries, err := os.ReadDir(wd)
	do.Fuck(err)

	relevantFiles := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			do.Verbose(entry.Name(), "is a directory")
			continue
		}
		if !isRelevantFile(entry.Name()) {
			do.Verbose(entry.Name(), "is not relevant")
			continue
		}
		do.Verbose(entry.Name(), "is added to the list")
		relevantFiles = append(relevantFiles, path.Join(wd, entry.Name()))
	}
	if do.Verbosity {
		for _, entry := range relevantFiles {
			do.Verbose(entry)
		}
	}
	game := NewGame(relevantFiles, wd)
	do.Fuck(ebiten.RunGame(game))
}
