package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	termbox "github.com/nsf/termbox-go"
	// termbox "github.com/nsf/termbox-go"
)

var (
	root = "."
	// debug    Debug
	// screen
)

//TODO: cache open files?
//TODO: slow version which opens and writes same file... kills SSDs...

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

// func addMatch(m Match) {
// 	gMatches = append(gMatches, m)
// }

// func xxxredraw(ev *termbox.Event) {
// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

// 	lastPath := ""
// 	_, h := termbox.Size()

// 	// Top line is for user input / status messages.
// 	y := 0
// 	for idx, m := range gMatches {
// 		x := 0

// 		// Print file name
// 		if m.path != lastPath {
// 			lastPath = m.path
// 			tbPrint(x, y, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault, m.path)
// 			y++
// 		}

// 		y = m.Print(x, y, idx == gSelected)

// 		// Dont draw off-screen
// 		if y > h {
// 			break
// 		}
// 	}

// 	// Vim style tildes for empty lines..
// 	for y < h-1 {
// 		tbPrint(0, y, termbox.ColorBlue|termbox.AttrBold, termbox.ColorDefault, "~")
// 		y++
// 	}

// 	// Dump debug info
// 	debugString := fmt.Sprintf("sel=%d", gSelected)
// 	tbPrint(0, h-2, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, debugString)

// 	// Status bar...
// 	tbPrint(0, h-1, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, "QUERY >>> ")
// 	tbPrint(10, h-1, termbox.ColorGreen|termbox.AttrBold|termbox.AttrReverse, termbox.ColorDefault, " ")

// 	termbox.Flush()
// }

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: qrr <pattern> <replacement>")
		os.Exit(-1)
	}
	regexFind := regexp.MustCompile(`\b` + os.Args[1] + `\b`)
	replaceWith := os.Args[2]

	screen := NewScreen()

	// meh?
	if _, ok := os.LookupEnv("TERMBOX256"); ok {
		termbox.SetOutputMode(termbox.Output256)
	}

	done := make(chan struct{})
	defer close(done)

	// TODO: check error channel
	matchesc, _ := processFiles(done, root, regexFind, replaceWith)
	eventsc := termboxPoll()

	// done <- struct{}{} // UGLY AF

mainloop:
	for {
		//debug.Add("failed")

		select {
		case ev := <-eventsc:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					break mainloop
				case termbox.KeyPgup:
					screen.selected = max(screen.selected-10, 0)
				case termbox.KeyPgdn:
					screen.selected = min(screen.selected+10, screen.totalMatchCount-1)
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
						break mainloop
					}
				}
			} else if ev.Type == termbox.EventResize {
				screen.w = ev.Width
				screen.h = ev.Height
			}

			screen.Redraw()

		case m, more := <-matchesc:
			if !more {
				// break Outer
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
		}
	}

	// time.Sleep(3 * time.Second)
	termbox.Close()
}
