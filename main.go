package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Counts struct {
	Bytes int
	Lines int
	Words int
	Chars int
}

func noArgs(args []string) bool {
	return len(args) < 2
}

func count(r io.Reader) (Counts, error) {
	reader := bufio.NewReader(r)
	var c Counts

	inWord := false

	for {
		rn, size, err := reader.ReadRune()

		if err == io.EOF {
			break
		}

		if err != nil {
			return c, err
		}

		c.Bytes += size
		c.Chars++

		if rn == '\n' {
			c.Lines++
		}

		if unicode.IsSpace(rn) {
			inWord = false
		} else {
			if !inWord {
				c.Words++
				inWord = true
			}
		}
	}

	return c, nil
}

func main() {
	if noArgs(os.Args) {
		fmt.Println("Enter a valid argument.")
		return
	}

	var (
		filename string
		reader   io.Reader
	)

	flags := make(map[string]bool)

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			flags[arg] = true
		} else {
			filename = arg
		}
	}

	if filename == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		reader = file
	}

	counts, err := count(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Default behaviour: lines, words, bytes
	if len(flags) == 0 {
		if filename == "" {
			fmt.Printf("%8d %8d %8d\n",
				counts.Lines,
				counts.Words,
				counts.Bytes)
		} else {
			fmt.Printf("%8d %8d %8d %s\n",
				counts.Lines,
				counts.Words,
				counts.Bytes,
				filename)
		}
		return
	}

	if flags["-c"] {
		if filename == "" {
			fmt.Printf("%8d\n", counts.Bytes)
		} else {
			fmt.Printf("%8d %s\n", counts.Bytes, filename)
		}
	}

	if flags["-l"] {
		if filename == "" {
			fmt.Printf("%8d\n", counts.Lines)
		} else {
			fmt.Printf("%8d %s\n", counts.Lines, filename)
		}
	}

	if flags["-w"] {
		if filename == "" {
			fmt.Printf("%8d\n", counts.Words)
		} else {
			fmt.Printf("%8d %s\n", counts.Words, filename)
		}
	}

	if flags["-m"] {
		if filename == "" {
			fmt.Printf("%8d\n", counts.Chars)
		} else {
			fmt.Printf("%8d %s\n", counts.Chars, filename)
		}
	}
}
