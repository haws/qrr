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
	debug           map[string]int
}

func NewScreen() Screen {
	screen := Screen{}

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}

	screen.rootFolder = "."
	screen.matches = make(map[string][]Match)
	screen.debug = make(map[string]int)
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

func (s *Screen) Debug(key string, val int) {
	s.debug[key] = val
}

func (s *Screen) Redraw() {
	termbox.Clear(defaultFgColor, defaultBgColor)

	// Top line is for user input / status messages.
	line := 0

	//  To iterate in alphabetical order.
	keys := []string{}
	maxFilePathWidth := 0

	for k := range s.matches {
		keys = append(keys, k)
		if len(k) > maxFilePathWidth {
			maxFilePathWidth = len(k)
		}
	}
	sort.Strings(keys)

	matchIdx := 0
	// How many matches fit?

	// Slow...?
	matchesCapacity := 0
Outer1:
	for _, filepath := range keys {
		line++
		for range s.matches[filepath] {
			line += 2
			matchesCapacity++
			if line >= s.height-2 {
				break Outer1
			}
		}
	}

	matchesSkip := s.selected - matchesCapacity

	s.Debug("line", line)
	s.Debug("matchesCapacity", matchesCapacity)
	s.Debug("matchesSkip", matchesSkip)
	line = 0

Outer:
	for _, filepath := range keys {
		f := s.matches[filepath]
		col := 0
		headerDrawn := false

		//TODO: only print this if the section below is non empty

		//summary := fmt.Sprintf("%sfilepath
		// Print this line if the "block" below is not empty.

		for _, m := range f {
			if matchIdx > matchesSkip {

				if !headerDrawn {
					tbPrint(col, line, defaultFilepathColor, defaultBgColor, filepath)
					summ := fmt.Sprintf("%d matches", len(f))
					tbPrint(maxFilePathWidth+1, line, defaultFilepathColor, defaultBgColor, summ)
					line++
					headerDrawn = true
				}

				line = m.Print(col, line, matchIdx == s.selected)
			}
			matchIdx++

			// Dont draw off-screen
			if line > s.height-3 {
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
	// tbPrint(s.width-20, s.height-1, defaultFgColor, defaultBgColor, debugString)
	//hiPrint(0, h-2, termbox.ColorGreen|termbox.AttrBold, "<	sel=%d>", s.selected)

	// Status bar...
	//tbPrint(0, s.h-1, defaultStatusColor, defaultBgColor, "QUERY >>> ")
	//x := hiPrint(0, s.height-1, defaultStatusColor, "Replace <%s> with <%s>? ", s.patternSearch, s.patternReplace)

	debugString := fmt.Sprintf("sel=%d ", s.selected)

	for k, v := range s.debug {
		debugString += fmt.Sprintf("%s=%d ", k, v)
	}

	tbPrint(0, s.height-1, defaultStatusColor, defaultBgColor, debugString)

	//TODO: do it like this?
	// s.UpdateStatus("Replace <%s> with <%s>? ", s.patternSearch, s.patternReplace)
	// s.PrintCursor(x, s.height-1)

	termbox.Flush()
}

func (s *Screen) replaceAllMatches(re *regexp.Regexp, repl string) {
	for _, filematches := range s.matches {
		for _, match := range filematches {
			match.Replace(re, repl)
		}
	}
}
