package main

import (
	"fmt"
	"regexp"
	"sort"

	termbox "github.com/nsf/termbox-go"
)

type EditBox int

type FilePath string

type Stats struct {
	MatchesFound  int
	FilesReplaced int
}

type Screen struct {
	cursorX         int
	cursorY         int
	activeEditBox   int
	edit            []EditBox
	patternSearch   *regexp.Regexp
	patternReplace  string
	matches         map[string][]Match
	totalMatchCount int
	selected        int // Selected match
	stats           Stats
}

func NewScreen() Screen {
	screen := Screen{}
	screen.matches = make(map[string][]Match)
	return screen
}

func (s *Screen) AddMatch(m Match) {
	s.matches[m.path] = append(s.matches[m.path], m)
	s.totalMatchCount++
}

func (s *Screen) PrintCursor(x, y int) {
	tbPrint(x, y, defaultCursorColor, defaultBgColor, " ")
}

func (s *Screen) Redraw() {
	termbox.Clear(defaultFgColor, defaultBgColor)

	_, h := termbox.Size()

	// Top line is for user input / status messages.
	line := 0
	matchIdx := 0

	//  To iterate in alphabetical order.
	keys := []string{}
	for k := range s.matches {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, filepath := range keys {
		f := s.matches[filepath]
		col := 0

		tbPrint(col, line, defaultFilepathColor, defaultBgColor, filepath)
		line++

		for _, m := range f {
			line = m.Print(col, line, matchIdx == s.selected)
			matchIdx++

			// Dont draw off-screen
			if line > h {
				break
			}
		}
	}

	// Vim style tildes for empty lines..
	for line < h-1 {
		tbPrint(0, line, defaultTildeColor, defaultBgColor, "~")
		line++
	}

	// Dump debug info
	debugString := fmt.Sprintf("sel=%d", s.selected)
	tbPrint(0, h-2, defaultFgColor, defaultBgColor, debugString)
	//hiPrint(0, h-2, termbox.ColorGreen|termbox.AttrBold, "<	sel=%d>", s.selected)

	// Status bar...
	tbPrint(0, h-1, defaultStatusColor, defaultBgColor, "QUERY >>> ")
	s.PrintCursor(10, h-1)

	termbox.Flush()
}

func (s *Screen) replaceAllMatches(re *regexp.Regexp, repl string) {
	for _, filematches := range s.matches {
		for _, match := range filematches {
			match.Replace(re, repl)
		}
	}
}
