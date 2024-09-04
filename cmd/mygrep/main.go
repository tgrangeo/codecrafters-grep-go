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
	group := []string{}
	result := strings.Builder{}
	i := 0

	for i < len(pattern) {
		if pattern[i] == '(' {
			j := i + 1
			depth := 1
			for j < len(pattern) && depth > 0 {
				if pattern[j] == '(' {
					depth++
				} else if pattern[j] == ')' {
					depth--
				}
				j++
			}
			if depth > 0 {
				return "", fmt.Errorf("unbalanced parentheses")
			}
			group = append(group, pattern[i:j])
			result.WriteString(pattern[i:j])
			i = j
		} else if pattern[i] == '\\' && i+1 < len(pattern) && pattern[i+1] >= '1' && pattern[i+1] <= '9' {
			refNum := int(pattern[i+1] - '0')
			if refNum <= len(group) {
				result.WriteString(group[refNum-1])
			} else {
				return "", fmt.Errorf("invalid back reference: \\%d", refNum)
			}
			i += 2
		} else {
			result.WriteByte(pattern[i])
			i++
		}
	}
	return result.String(), nil
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
