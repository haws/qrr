package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	termbox "github.com/nsf/termbox-go"
	// termbox "github.com/nsf/termbox-go"
)

var (
	filesProcessed int
)

//TODO: cache open files?
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
					filesProcessed++
					scanner := bufio.NewScanner(f)
					lineNo := 1
					for scanner.Scan() {
						lineFrom := strings.TrimSpace(scanner.Text())
						//gMatches := reFrom.FindAllString(lineFrom, -1)

						linematches := reFrom.FindAllStringIndex(lineFrom, -1)

						if linematches != nil {
							newline := reFrom.ReplaceAllString(lineFrom, replaceWith)
							matchc <- Match{
								lineNo:      lineNo,
								path:        path,
								line:        lineFrom,
								newline:     newline,
								linematches: linematches,
								repl:        replaceWith,
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

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: qrr <pattern> <replacement>")
		os.Exit(-1)
	}
	regexFind := regexp.MustCompile(`\b` + regexp.QuoteMeta(os.Args[1]) + `\b`)
	replaceWith := os.Args[2]

	screen := NewScreen()
	screen.patternSearch = os.Args[1]
	screen.patternReplace = os.Args[2]

	// meh?
	if _, ok := os.LookupEnv("TERMBOX256"); ok {
		termbox.SetOutputMode(termbox.Output256)
	}

	done := make(chan struct{})
	defer close(done)

	// TODO: check error channel
	matchesc, _ := processFiles(done, screen.rootFolder, regexFind, replaceWith)
	eventsc := termboxPoll()
	ticker := time.NewTicker(time.Millisecond * 500)
	// done <- struct{}{} // UGLY AF

mainloop:
	for {
		//debug.Add("failed")

		select {
		case ev := <-eventsc:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					screen.Close()
					fmt.Println("Quit.")
					break mainloop
				case termbox.KeyPgup:
					matchesInWindow := (screen.height - 1)
					screen.selected = max(screen.selected-matchesInWindow-1, 0)
				case termbox.KeyPgdn:
					matchesInWindow := (screen.height - 1)
					screen.selected = min(screen.selected+matchesInWindow-1, screen.totalMatchCount-1)
				case termbox.KeyHome:
					screen.selected = 0
				case termbox.KeyEnd:
					screen.selected = screen.totalMatchCount - 1
				case termbox.KeyArrowUp:
					screen.selected = max(screen.selected-1, 0)
				case termbox.KeyArrowDown:
					screen.selected = min(screen.selected+1, screen.totalMatchCount-1)
				case termbox.KeyEnter:
					screen.replaceAllMatches(regexFind, replaceWith)
					screen.Close()
					screen.PrintStats()
					break mainloop
					// TODO: replace and jump
					//gSelected = min(gSelected+1, len(gMatches)-1)

				// case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				// 	edit_box.MoveCursorOneRuneBackward()
				// case termbox.KeyArrowRight, termbox.KeyCtrlF:
				// 	edit_box.MoveCursorOneRuneForward()
				// case termbox.KeyBackspace, termbox.KeyBackspace2:
				// 	edit_box.DeleteRuneBackward()
				// case termbox.KeyDelete, termbox.KeyCtrlD:
				// 	edit_box.DeleteRuneForward()
				// case termbox.KeyTab:
				// 	edit_box.InsertRune('\t')
				// case termbox.KeySpace:
				// 	edit_box.InsertRune(' ')
				// case termbox.KeyCtrlK:
				// 	edit_box.DeleteTheRestOfTheLine()
				// case termbox.KeyHome, termbox.KeyCtrlA:
				// 	edit_box.MoveCursorToBeginningOfTheLine()
				// case termbox.KeyEnd, termbox.KeyCtrlE:
				// 	edit_box.MoveCursorToEndOfTheLine()
				default:
					if ev.Ch == 'q' || ev.Ch == 'Q' {
						screen.Close()
						fmt.Println("Quit.")
						break mainloop
					}
				}
			} else if ev.Type == termbox.EventResize {
				screen.width = ev.Width
				screen.height = ev.Height
			}

			screen.Redraw()

		case m, more := <-matchesc:
			if !more {
				if screen.state != stateDone {
					screen.Done()
					screen.Redraw()
				}
			} else {
				screen.AddMatch(m)
				screen.Redraw()
				//addMatch(m)
				// redraw(nil)
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
		case <-ticker.C:
			screen.Redraw()
		}
	}

	screen.Close()

}
