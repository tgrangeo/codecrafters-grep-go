package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
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

func checkBackReferences(pattern string) (string, error) {
	var groups []string
	var newPattern strings.Builder
	groupCounter := 1

	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '(':
			start := i
			end := i + strings.Index(pattern[i:], ")")
			groups = append(groups, pattern[start:end+1])
			groupCounter++
			newPattern.WriteString(pattern[start:end+1])
			i = end
		case '\\':
			if i+1 < len(pattern) && pattern[i+1] >= '1' && pattern[i+1] <= '9' {
				refNum := int(pattern[i+1] - '0')
				if refNum <= len(groups) {
					newPattern.WriteString(groups[refNum-1])
				} else {
					return "", fmt.Errorf("invalid back reference: \\%d", refNum)
				}
				i++
			} else {
				newPattern.WriteByte(pattern[i])
			}
		default:
			newPattern.WriteByte(pattern[i])
		}
	}

	return newPattern.String(), nil
}

func matchLine(line []byte, pattern string) (bool, error) {
	pattern, err := checkBackReferences(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid back reference: %v", err)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid pattern: %v", err)
	}
	return re.Match(line), nil
}
