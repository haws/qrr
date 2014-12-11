package main 

import (
	"os"
	"fmt"
	"regexp"
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

func redraw(w int, h int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	tbprint(0,h-1,termbox.ColorRed, termbox.ColorWhite, "  hello            ")
	termbox.SetCell(0,0,'a', termbox.ColorRed|termbox.AttrUnderline, termbox.ColorDefault)
	termbox.SetCell(20,20,'b', termbox.ColorBlue|termbox.AttrReverse, termbox.ColorDefault)
	termbox.SetCell(w-1,h-1,'c', termbox.ColorGreen, termbox.ColorDefault)
	termbox.Flush()
}

func try_replace(path string, from *regexp.Regexp, to *regexp.Regexp) {
	termbox.Flush()
	a, b := termbox.Size()
	redraw(a,b)

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventResize {
			redraw(ev.Width, ev.Height)
		} else {
			termbox.Close()
			os.Exit(0)
		}
	}
}

func main() {
      usage := `query replace regexp 0.0.1

Search for pattern <from> in each file in the current directory, asking the
user if he wants to replace it with <to>. The proposed change is displayed as
a colored diff.

The patterns are Perl-style regexes. 

Usage:
  qrr <from> <to> [-b] [-p <project>]
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
	re_from := regexp.MustCompile(args["<from>"].(string))
	re_to   := regexp.MustCompile(args["<to>"].(string))

	err := termbox.Init()
	defer termbox.Close()

	if err != nil {
		fmt.Println("Could not start termbox. View ~/.codechange.log for error messages.")
		os.Exit(1)
	}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			try_replace(path, re_from, re_to)
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
