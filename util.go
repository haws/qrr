package main

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

var (
	defaultBgColor       = termbox.ColorDefault
	defaultFgColor       = termbox.ColorDefault
	defaultCursorColor   = termbox.ColorGreen | termbox.AttrBold | termbox.AttrReverse
	defaultTildeColor    = termbox.ColorGreen | termbox.AttrBold
	defaultFilepathColor = termbox.ColorCyan | termbox.AttrBold
	defaultLineColor     = termbox.ColorYellow | termbox.AttrBold
	defaultRemovedColor  = termbox.ColorYellow | termbox.AttrBold
	defaultAddedColor    = termbox.ColorGreen | termbox.AttrBold
	defaultStatusColor   = termbox.ColorRed | termbox.AttrBold
)

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func hiPrint(x, y int, hi termbox.Attribute, format string, a ...interface{}) {
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault

	curfg := fg

	s := fmt.Sprintf(format, a...)
	for _, c := range s {
		if c == '<' {
			curfg = hi
		} else if c == '>' {
			curfg = fg
		} else {
			termbox.SetCell(x, y, c, curfg, bg)
			x++
		}
	}

}

// type Debug struct {
// 	msgs []string
// }

// func (d Debug) Add(format string, a ...interface{}) {
// 	s := fmt.Sprintf(format, a...)
// 	d.msgs = append(d.msgs, s)
// }

// func (d Debug) Print() {
// 	_, h := termbox.Size()
// 	x := 0
// 	tbPrint(x, h-1, termbox.ColorDefault, termbox.ColorDefault, fmt.Sprintf("%d", len(d.msgs)))
// 	// for _, s := range d.msgs {
// 	// 	tbPrint(x, h-1, termbox.ColorDefault, termbox.ColorDefault, s)
// 	// 	x += len(s)
// 	// }
// 	// d.msgs = make([]string, 1)
// }

// func debugPrint(format string, a ...interface{}) {
// 	s := fmt.Sprintf(format, a...)
// 	_, h := termbox.Size()
// 	tbPrint(0, h-1, termbox.ColorRed, termbox.ColorDefault, s)
// }

// func (d Debug) Print(format string, a ...interface{}) {
// 	s := fmt.Sprintf(format, a...)
// 	_, h := termbox.Size()
// 	tbPrint(0, h-1, termbox.ColorDefault, termbox.ColorDefault, s)
// }

func termboxPoll() chan termbox.Event {
	evCh := make(chan termbox.Event)

	go func() {
		for {
			evCh <- termbox.PollEvent()
		}
	}()

	return evCh
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
