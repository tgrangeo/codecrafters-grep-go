package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	// "unicode/utf8"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	// if utf8.RuneCountInString(pattern) != 1 {
	// 	return false, fmt.Errorf("unsupported pattern: %q", pattern)
	// }
	if bytes.ContainsAny(line, pattern) {
		return true, nil
	} else if pattern == `\d` {
		if bytes.ContainsAny(line, "0123456789") {
			return true, nil
		}
	} else if pattern == `\w` {
		re := regexp.MustCompile(`[a-zA-Z0-9]`)
		if re.Match(line) {
			return true, nil
		}
	} else if regexp.MustCompile(`^\[\^[a-zA-Z0-9]+\]$`).MatchString(pattern){
		re := regexp.MustCompile(`^\[\^[a-zA-Z0-9]+\]$`)
		matches := re.FindSubmatch([]byte(pattern))
		if len(matches) > 1 {
			toFind :=  string(matches[1])
			if !bytes.ContainsAny(line, toFind) {
				return true, nil
			}
		}
	}
	return false, nil
}
