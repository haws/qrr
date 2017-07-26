package main

import (
	"fmt"
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

func (m Match) Print(initialX, initialY int, isSelected bool) int {
	x, y := initialX, initialY

	lineColor := defaultLineColor
	fgColor := defaultFgColor
	bgColor := defaultBgColor
	removedColor := defaultRemovedColor
	addedColor := defaultAddedColor

	if isSelected {
		fgColor = fgColor | termbox.AttrReverse
		bgColor = bgColor | termbox.AttrReverse
		removedColor = removedColor | termbox.AttrReverse
		addedColor = addedColor | termbox.AttrReverse
	}

	// First line
	lineNumber := fmt.Sprintf("%4d\t", m.lineNo)

	tbPrint(x, y, lineColor, termbox.ColorDefault, lineNumber)

	for _, sm := range m.linematches {
		beg := sm[0]
		end := sm[1]
		tbPrint(x+len(lineNumber), y, fgColor, bgColor, m.line[x:beg])
		tbPrint(beg+len(lineNumber), y, removedColor, bgColor, m.line[beg:end])
		x = end
	}
	tbPrint(x+len(lineNumber), y, fgColor, bgColor, m.line[x:])

	// Second line
	x = initialX
	y++
	origStringIdx := 0
	// w, _ := termbox.Size()
	xoff := 0
	// tbPrint(x, y, termbox.ColorGreen|termbox.AttrBold, bgColor, lineNumber)

	for _, sm := range m.linematches {
		beg := sm[0]
		end := sm[1]
		tbPrint(xoff+x+len(lineNumber), y, fgColor, bgColor, m.line[origStringIdx:beg])
		x += (beg - origStringIdx)
		tbPrint(xoff+x+len(lineNumber), y, addedColor, bgColor, m.repl)
		x += len(m.repl)
		origStringIdx = end
	}
	tbPrint(xoff+x+len(lineNumber), y, fgColor, bgColor, m.line[origStringIdx:])

	return y + 1
}
