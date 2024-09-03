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

func checkBackReferences(line []byte, pattern string) (string, error) {
	if strings.Contains(pattern, "\\1") {
		res := pattern
		begin := strings.IndexRune(res, '(')
		end := strings.IndexRune(res, ')')
		group := res[begin : end+1]
		// res = strings.ReplaceAll(res, "(", "")
		// res = strings.ReplaceAll(res, ")", "")
		res = strings.ReplaceAll(res, `\1`, group)
		return res, nil
	}
	return pattern, nil
}

func matchLine(line []byte, pattern string) (bool, error) {
	pattern, err := checkBackReferences(line, pattern)
	fmt.Println(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid back reference: %v", err)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid pattern: %v", err)
	}
	return re.Match(line), nil
}
