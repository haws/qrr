package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	termbox "github.com/nsf/termbox-go"
	// termbox "github.com/nsf/termbox-go"
)

type Match struct {
	path    string // Filepath
	lineNo  int    // Line number
	line    string // Line with maches
	newline string // Line with replacements
}

var (
	root    = "."
	matches []Match
	// screen
)

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		// Close the paths channel after Walk returns.
		defer close(paths)
		// No select needed for this send, since errc is buffered.
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && ignoredFolders[info.Name()] {
				return filepath.SkipDir
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			if !allowedExtensions[filepath.Ext(path)] {
				return nil
			}

			select {
			case paths <- path:
			case <-done:
				return errors.New("walk canceled")
			}
			return nil
		})

	}()
	return paths, errc
}

func processFiles(done <-chan struct{}, root string, reFrom *regexp.Regexp, replaceWith string) (<-chan Match, <-chan error) {
	paths, errc := walkFiles(done, root)
	matchc := make(chan Match)

	go func() {
		defer close(matchc)

		for {
			select {
			case path, more := <-paths:
				if !more {
					return
				}
				f, err := os.Open(path)
				if err == nil {
					scanner := bufio.NewScanner(f)
					lineNo := 1
					for scanner.Scan() {
						lineFrom := scanner.Text()
						matches := reFrom.FindAllString(lineFrom, -1)

						if matches != nil {
							newline := reFrom.ReplaceAllString(lineFrom, replaceWith)
							matchc <- Match{
								lineNo:  lineNo,
								path:    path,
								line:    lineFrom,
								newline: newline,
							}
						}
						lineNo++
					}
				} else {
					fmt.Println(err)
				}
			case <-done:
				fmt.Println("got a done")
				return
			}
		}

	}()

	return matchc, errc
}

func addMatch(m Match) {
	matches = append(matches, m)
}

func redraw(ev *termbox.Event) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	lastPath := ""

	_, h := termbox.Size()

	y := 1
	for _, m := range matches {
		x := 0
		if m.path != lastPath {
			lastPath = m.path
			tbPrint(x, y, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault, m.path)
			y++
		}
		if y < h {
			tbPrint(x, y, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault, fmt.Sprintf("%4d  ", m.lineNo))
			x += 6
			tbPrint(x, y, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault, fmt.Sprintf("%s", m.line))

			x = 0
			tbPrint(x, y+1, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault, fmt.Sprintf("%4d  ", m.lineNo))
			x += 6
			tbPrint(x, y+1, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, fmt.Sprintf("%s", m.newline))
			y += 2
		}
	}

	termbox.Flush()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: qrr <pattern> <replacement>")
		os.Exit(-1)
	}
	regexFind := regexp.MustCompile(`\b` + os.Args[1] + `\b`)
	replaceWith := os.Args[2]

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	defer close(done)

	// TODO: check error channel
	matchesc, _ := processFiles(done, root, regexFind, replaceWith)
	eventsc := termboxPoll()

	// done <- struct{}{} // UGLY AF

Outer:
	for {
		select {
		case ev := <-eventsc:
			if ev.Type == termbox.EventKey { //&& ev.Key == termbox.KeyEsc {
				break Outer
			}
			redraw(&ev)

		case m, more := <-matchesc:
			if !more {
				// break Outer
			} else {
				addMatch(m)
				redraw(nil)
			}
			// fmt.Printf("- %s:%d %s\n", m.path, m.lineNo, m.line)
			// fmt.Printf("+ %s:%d %s\n", m.path, m.lineNo, m.newline)
			//fmt.Println("match", m.lineNo, m.path, m.line)
			// fmt.Println(m.path)
			// fmt.Println(m.line)
			// fmt.Println(m.newline)

			// case <-done:
			// fmt.Println("finish")
			// break Outer
			// case err := <-errc:

			// fmt.Println("read from errc", err)
			// log.Fatal(err)
		}
	}

	// time.Sleep(3 * time.Second)
	termbox.Close()
}
