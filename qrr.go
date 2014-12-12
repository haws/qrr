package main 

import (
	"os"
	"fmt"
	"bytes"
	"regexp"
	"io/ioutil"
	"path/filepath"

	"github.com/nsf/termbox-go"
	"github.com/docopt/docopt-go"
)

var (
	ignored_dirs = map[string]bool {
		".git": true,
	}
	root = "."
)

// Information about a regex match in the buffer
type MatchLine struct {
    startIndex  int 	// Buffer position where match line starts
    startLine   int 	// Line where match starts
    endIndex    int 	// Buffer position where match line ends
    endLine     int 	// Line where match ends (startLine+1 if match on a single line, etc)
}


func findNewlinesIndex(buf []byte) []int {
	result := make([]int, 0, 10)
	result = append(result,0) // first "newline" in 0 position simplifies things

	i := 0
	offset := 0

	for {
		if i = bytes.IndexByte(buf, '\n'); i >= 0 {
			buf = buf[i+1:]
			offset = offset + i + 1
			result = append(result, offset-1)
		} else {
			break
		}
	}

	return result
}

func exit() {
	termbox.Close()
	os.Exit(2)
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func redraw(w int, h int, s string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	tbprint(0,h-1,termbox.ColorRed, termbox.ColorWhite, "  hello            ")
	termbox.SetCell(0,0,'a', termbox.ColorRed|termbox.AttrUnderline, termbox.ColorDefault)
	termbox.SetCell(20,20,'b', termbox.ColorBlue|termbox.AttrReverse, termbox.ColorDefault)
	termbox.SetCell(w-1,h-1,'c', termbox.ColorGreen, termbox.ColorDefault)
	tbprint(5,5,termbox.ColorYellow, termbox.ColorDefault,s)
	termbox.Flush()
}

func getMatchLine(b []byte, newlines []int, match_begin int, match_end int) (m MatchLine) {
	for line_count, newline_pos := range newlines {
		if match_begin >= newline_pos {
			if newline_pos > 0 {
				m.startIndex = newline_pos+1
			}
			m.startLine = line_count+1
		}
	}
	for line_count, newline_pos := range newlines {
		if match_end <= newline_pos {
			m.endIndex = newline_pos
			m.endLine = line_count
			break
		}
	}

	return m
}

func try_replace(path string, pattern *regexp.Regexp, replace string) {
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	matches := pattern.FindAllIndex(contents, -1)
	newlines := findNewlinesIndex(contents)

	for _, match := range matches {
		//fmt.Println(newlines)
		ml := getMatchLine(contents, newlines, match[0], match[1])
		//fmt.Println(path, newlines, string(contents[match[0]:match[1]]))
		fmt.Printf("%3d %12s %s\n", ml.startLine, path, string(contents[ml.startIndex:ml.endIndex]))
	}

	//fmt.Println(path, matches, newlines)

	//fmt.Println(len(findNewlinesIndex(contents)))
/*
	termbox.Flush()
	a, b := termbox.Size()
	
	redraw(a,b,path)

	ev := termbox.PollEvent()
	if ev.Type == termbox.EventResize {
		redraw(ev.Width, ev.Height,path)
	}
*/
}

func main() {
      usage := `query replace regexp 0.0.1

Search for pattern <pattern> in each file in the current directory, asking the
user if he wants to replace it with <new>. The proposed change is displayed as
a colored diff.

The patterns are Perl-style regexes. 

Usage:
  qrr <pattern> <new> [-b] [-p <project>]
  qrr --genconfig
  qrr -h | --help
  qrr --version

Options:
  -h --help                        Show this screen.
  --version                        Show version.
  --genconfig                      Generate example config file.
  -p <project> --project=<project> Replace in <project> instead of directory (see --genconfig).
  -b --bulk                        Bulk mode (quicker, but doesn't shows context lines).
`
	args, _ := docopt.Parse(usage, nil, true, "qrr 0.0.1", false)
	pattern := regexp.MustCompile(args["<pattern>"].(string))
	replace := args["<new>"].(string)

	//err := termbox.Init()
	//defer termbox.Close()

	//if err != nil {
	//	fmt.Println("Could not start termbox. View ~/.codechange.log for error messages.")
	//	os.Exit(1)
	//}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && path == "README.md" {
			try_replace(path, pattern, replace)
		} else {
			if ignored_dirs[info.Name()] {
				return filepath.SkipDir
			}
		}
		return nil
	})

}



// .\codemod.py:70


// - You can also use codemod for transformations that are much more sophisticated
// + You can also use codemod for transformations that are much more sofisticated
//   than regular expression substitution.  Rather than using the command line, you

// Accept change (y = yes [default], n = no, e = edit, E = yes+edit, q = quit)?
