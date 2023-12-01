package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"strconv"
	"unicode"
)

//go:embed input.txt
var input embed.FS

func main() {
	b, err := input.ReadFile("input.txt")

	if err != nil {
		log.Fatalf("failed to open input: %s", err)
	}

	reader := bytes.NewReader(b)
	buf := bufio.NewReader(reader)

	var sum int64

	for {
		line, err := buf.ReadBytes('\n')
		lineBuf := bufio.NewReader(bytes.NewReader(line))

		var first *rune
		var last *rune
		for {
			r, _, err := lineBuf.ReadRune()
			if err == io.EOF {
				// No more bytes to read on this line
				break
			} else if err != nil {
				log.Fatalf("Got error reading rune from line: %s", err)
			}

			log.Printf("Examining rune %s", string(r))

			if unicode.IsNumber(r) {
				if first == nil {
					log.Printf("Got first int %s", string(r))
					first = &r
				}

				log.Printf("Also updating last int %s", string(r))
				// TODO I guess this one doesn't need to be a pointer after all
				last = &r
			}
		}

		concat, parseErr := strconv.ParseInt(fmt.Sprintf("%s%s", string(*first), string(*last)), 10, 64)
		if parseErr != nil {
			log.Fatalf("failed to mush integers together: %s", parseErr)
		}

		log.Printf("Got first, last: %s, %s. Concatenated that makes: %d", string(*first), string(*last), concat)

		sum = sum + concat

		if err == io.EOF {
			log.Println("Reached end of file")
			break
		}
	}

	log.Printf("Got sum total: %d", sum)
}
