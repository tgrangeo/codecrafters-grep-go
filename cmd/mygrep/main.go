package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin)
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
}

func getPattern(patternLine string) []string {
	res := []string{}
	for i := 0; i < len(patternLine); i++ {
		if patternLine[i] == '\\' && i+1 < len(patternLine) {
			res = append(res, patternLine[i:i+2])
			i++
		} else {
			res = append(res, string(patternLine[i]))
		}
	}
	return res
}

func checkNextPattern(c byte, p string) bool {
	switch p {
	case `\d`:
		return bytes.ContainsAny([]byte{c}, "0123456789")
	case `\w`:
		return regexp.MustCompile(`[a-zA-Z0-9]`).Match([]byte{c})
	default:
		if p[0] == '[' && p[len(p)-1] == ']' {
			if p[1] == '^' {
				return !bytes.ContainsAny([]byte{c}, p[2:len(p)-1])
			}
			return bytes.ContainsAny([]byte{c}, p[1:len(p)-1])
		}
		return string(c) == p
	}
}

func matchLine(line []byte, pattern string) (bool, error) {
	patternArray := getPattern(pattern)
	lineIndex := 0

	for _, p := range patternArray {
		if lineIndex >= len(line) {
			return false, nil
		}

		if !checkNextPattern(line[lineIndex], p) {
			return false, nil
		}

		lineIndex++
	}

	return true, nil
}
