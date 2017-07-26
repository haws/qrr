package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

type Match struct {
	path        string  // Filepath
	lineNo      int     // Line number
	line        string  // Line with maches
	newline     string  // Line with replacements
	repl        string  // Replacement string
	linematches [][]int // Positions of gMatches
	marked      bool    // Line should be replaced? (TODO: How to replace only a few in a line)
}

func (m Match) Replace(re *regexp.Regexp, repl string) {
	input, err := ioutil.ReadFile(m.path)
	if err != nil {
		termbox.Close()
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	// Replacement
	lines[m.lineNo-1] = re.ReplaceAllString(lines[m.lineNo-1], repl)

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(m.path, []byte(output), 0644)
	if err != nil {
		termbox.Close()
		log.Fatalln(err)
	}

	// if m.marked {
	// }
}
