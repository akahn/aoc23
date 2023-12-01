package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"unicode"
)

//go:embed input.txt
var input embed.FS

func p1() {
	b, err := input.ReadFile("input.txt")

	if err != nil {
		log.Fatalf("failed to open input: %s", err)
	}

	reader := bytes.NewReader(b)
	buf := bufio.NewReader(reader)

	var sum int

	for {
		line, err := buf.ReadBytes('\n')
		lineBuf := bufio.NewReader(bytes.NewReader(line))
		number := extractNumber(lineBuf)

		sum = sum + number

		if err == io.EOF {
			log.Println("Reached end of file")
			break
		}
	}

	log.Printf("Got part 1 sum total: %d", sum)
}

func p2() {
	b, err := input.ReadFile("input.txt")

	if err != nil {
		log.Fatalf("failed to open input: %s", err)
	}

	reader := bytes.NewReader(b)
	buf := bufio.NewReader(reader)

	var sum int
	for {
		line, err := buf.ReadBytes('\n')
		lineBuf := bufio.NewReader(bytes.NewReader(line))
		number := extractNumberIncludingFromWords(lineBuf)

		sum = sum + number

		if err == io.EOF {
			log.Println("Reached end of file")
			break
		}
	}

	log.Printf("Got part 2 sum total: %d", sum)
}
func main() {
	p1()
	p2()
}

var detectNumbers = regexp.MustCompile("(\\d|one|two|three|four|five|six|seven|eight|nine)")
var table = map[string]int{
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func extractNumberIncludingFromWords(line io.Reader) int {
	b, err := io.ReadAll(line)
	log.Printf("Going to regex %s", string(b))
	if err != nil {
		log.Fatalf("Failed to read input line: %s", err)
	}
	matches := detectNumbers.FindAllStringSubmatch(string(b), -1)
	length := len(matches)
	firstMatch := matches[0][0]
	lastMatch := matches[length-1][0]

	firstInt, ok := table[firstMatch]
	if !ok {
		log.Fatalf("Failed to look up first match in table %s", firstMatch)
	}
	lastInt, ok := table[lastMatch]
	if !ok {
		log.Fatalf("Failed to look up last match in table %s", lastMatch)
	}

	concat, err := strconv.ParseInt(fmt.Sprintf("%d%d", firstInt, lastInt), 10, 64)
	if err != nil {
		log.Fatalf("failed to mush integers together: %s", err)
	}
	log.Printf("Got %d matches. First, last: %s, %s. Concat: %d", length, firstMatch, lastMatch, concat)

	return int(concat)
}

func extractNumber(line io.Reader) int {
	lineBuf := bufio.NewReader(line)

	var first *rune
	var last rune
	for {
		r, _, err := lineBuf.ReadRune()
		if err == io.EOF {
			// No more bytes to read on this line
			break
		} else if err != nil {
			log.Fatalf("Got error reading rune from line: %s", err)
		}

		//log.Printf("Examining rune %s", string(r))

		if unicode.IsNumber(r) {
			if first == nil {
				//log.Printf("Got first int %s", string(r))
				first = &r
			}

			//log.Printf("Also updating last int %s", string(r))
			last = r
		}
	}

	concat, err := strconv.ParseInt(fmt.Sprintf("%s%s", string(*first), string(last)), 10, 64)
	if err != nil {
		log.Fatalf("failed to mush integers together: %s", err)
	}

	return int(concat)
}
