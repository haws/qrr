package main

import (
	"fmt"
	"log"
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
	line, col       int
	width           int // width of the screen
	height          int // height of the screen
	activeEditBox   int
	edit            []EditBox
	rootFolder      string
	patternSearch   string
	patternReplace  string
	matches         map[string][]Match
	totalMatchCount int
	selected        int // Selected match
	stats           Stats
}

func NewScreen() Screen {
	screen := Screen{}

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}

	screen.rootFolder = "."
	screen.matches = make(map[string][]Match)
	screen.width, screen.height = termbox.Size()
	return screen
}

func (s *Screen) NextLine() {
	s.line++
	s.col = 0
}

func (s *Screen) Print(fg, bg termbox.Attribute, format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)

	for _, c := range str {
		termbox.SetCell(s.col, s.line, c, fg, bg)
		s.col++
	}
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

	// Top line is for user input / status messages.
	line := 0
	matchIdx := 0

	//  To iterate in alphabetical order.
	keys := []string{}
	for k := range s.matches {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// How many matches fit?
	matchesCapacity := (s.height - 1) / 2
	matchesSkip := s.selected - matchesCapacity

Outer:
	for _, filepath := range keys {
		f := s.matches[filepath]
		col := 0

		tbPrint(col, line, defaultFilepathColor, defaultBgColor, filepath)
		line++

		for _, m := range f {
			if matchIdx > matchesSkip {
				line = m.Print(col, line, matchIdx == s.selected)
			}
			matchIdx++

			// Dont draw off-screen
			if line > s.height-2 {
				break Outer
			}
		}
	}

	// Vim style tildes for empty lines..
	for line < s.height-1 {
		tbPrint(0, line, defaultTildeColor, defaultBgColor, "~")
		line++
	}

	// Dump debug info
	debugString := fmt.Sprintf("sel=%d", s.selected)
	tbPrint(s.width-20, s.height-1, defaultFgColor, defaultBgColor, debugString)
	//hiPrint(0, h-2, termbox.ColorGreen|termbox.AttrBold, "<	sel=%d>", s.selected)

	// Status bar...
	//tbPrint(0, s.h-1, defaultStatusColor, defaultBgColor, "QUERY >>> ")
	x := hiPrint(0, s.height-1, defaultStatusColor, "Replace <%s> with <%s>? ", s.patternSearch, s.patternReplace)

	//TODO: do it like this?
	// s.UpdateStatus("Replace <%s> with <%s>? ", s.patternSearch, s.patternReplace)
	s.PrintCursor(x, s.height-1)

	termbox.Flush()
}

func (s *Screen) replaceAllMatches(re *regexp.Regexp, repl string) {
	for _, filematches := range s.matches {
		for _, match := range filematches {
			match.Replace(re, repl)
		}
	}
}
