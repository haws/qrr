package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	flags "github.com/jessevdk/go-flags"
	termbox "github.com/nsf/termbox-go"
)

// Information about a regex match in the buffer
type MatchLine struct {
	startIndex int // Buffer position where match line starts
	startLine  int // Line where match starts
	endIndex   int // Buffer position where match line ends
	endLine    int // Line where match ends (startLine+1 if match on a single line, etc)
}

var opts struct {
	Recursive []bool `short:"r" long:"recursive" description:"Find files recursively"`
}

var (
	root = "."
)

func processFile(path string, regexFind *regexp.Regexp) {
	namedPrint := false
	f, err := os.Open(path)
	if err == nil {
		scanner := bufio.NewScanner(f)
		line := 1
		for scanner.Scan() {
			matches := regexFind.FindAllString(scanner.Text(), -1)
			if matches != nil {
				if !namedPrint {
					fmt.Printf("--- %s ----------------------------------------------------------------------------------------\n", path)
					namedPrint = true
				}
				newline := regexFind.ReplaceAllString(scanner.Text(), "newpattern")
				fmt.Println(scanner.Text())
				fmt.Println(newline)
			}

			// if strings.Contains(scanner.Text(), regexFind) {
			// 	fmt.Println(line, scanner.Text())
			// 	//return line, nil
			// }

			line++
		}
		// return 0, err
	}
	//fmt.Println(path, regexFind)
}

func main() {

	args, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	if len(args) != 1 {
		fmt.Println("usage: qrr <pattern>")
		os.Exit(-1)
	}
	regexFind := regexp.MustCompile(args[0])

	err = termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	termbox.Close()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if allowedExtensions[filepath.Ext(path)] {
				processFile(path, regexFind)
			}
		} else {
			if ignoredFolders[info.Name()] {
				return filepath.SkipDir
			}
		}
		return nil
	})

}
