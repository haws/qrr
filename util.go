package main

import (
	termbox "github.com/nsf/termbox-go"
)

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
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
