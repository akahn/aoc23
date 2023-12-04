package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

//go:embed input.txt
var input embed.FS

const placeHolder = "‽"

type Location struct {
	char   string
	number *string
}

type Schematic struct {
	grid [140][140]Location
	// Slice of points where parts have been found
	parts [][]int
}

func (s *Schematic) Set(x, y int, value string) {
	//log.Printf("Setting point %d,%d to %s", x, y, value)
	s.grid[y][x] = Location{char: value}

	if value != "." && value != placeHolder {
		//log.Printf("Stashing part %s at (%d,%d)", value, x, y)
		s.parts = append(s.parts, []int{x, y})
	}
}

func (s *Schematic) SetNumber(x, y int, value string) {
	// Set each cell to the full value, because that will make life easier later ???🤷???

	// Backtrack through the last n cells
	length := len(value)
	for i := 0; i < length; i++ {
		// Move left 1 due to starting from the next symbol after the number is finished
		// Then move further by the backtracking amount
		xPos := x - 1 - i
		if xPos < 0 {
			log.Fatalf("Error trying to set %s into the grid", value)
		}
		char := (value)[length-(i+1)]
		location := Location{char: string(char), number: &value}
		s.grid[y][xPos] = location
	}
}

func (s *Schematic) Neighbors(x, y int) []Location {
	neighbors := []Location{}
	// From x-1 to x+1 and y-1 to y+1 without coloring outside the lines

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			nX := x + i
			nY := y + j

			if nX < 0 || nY < 0 || nX > 139 || nY > 139 ||
				// No need to consider self
				(nX == 0 && nY == 0) {
				continue
			}

			neighbors = append(neighbors, s.grid[nY][nX])
		}
	}

	return neighbors
}

func (s *Schematic) Get(x, y int) string {
	return s.grid[y][x].char
}

func (s *Schematic) String() string {
	sb := strings.Builder{}

	for _, row := range s.grid {
		for _, location := range row {
			sb.WriteString(location.char)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (s *Schematic) HTML() string {
	sb := strings.Builder{}

	sb.WriteString("<code>\n")
	for y, row := range s.grid {
		for x, location := range row {
			hasPartNumber := false

			if location.char != "." && location.number == nil {
				for _, n := range s.Neighbors(x, y) {
					if n.number != nil {
						hasPartNumber = true
						break
					}

				}
			}
			color := "black"
			if hasPartNumber {
				color = "red"
			}
			sb.WriteString(fmt.Sprintf("<strong style='color:%s'>%s</strong>", color, location.char))
		}
		sb.WriteString("<br>\n")
	}

	sb.WriteString("</code>")
	return sb.String()
}

func parseSchematic() Schematic {
	b, err := input.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("failed to open input: %s", err)
	}

	reader := bytes.NewReader(b)

	schematic := Schematic{grid: [140][140]Location{}, parts: make([][]int, 0)}

	x := 0
	y := 0
	currentNumber := ""
	inNumber := false
	for {
		b, err := reader.ReadByte()

		if err != nil {
			log.Println("Reached end of file")
			break
		}

		if unicode.IsDigit(rune(b)) {
			inNumber = true
			currentNumber = currentNumber + string(b)
			schematic.Set(x, y, placeHolder)
		} else {
			if inNumber {
				// Encountered a non-number while in a number. Clear the number
				schematic.SetNumber(x, y, currentNumber)
				currentNumber = ""
				inNumber = false
			}

			// "Carriage return"
			if b == '\n' {
				y = y + 1
				x = 0
				continue
			}
			// TODO normalize symbols
			schematic.Set(x, y, string(b))
		}

		x = x + 1
	}

	return schematic
}

func p1() int {
	schematic := parseSchematic()
	partNeighbors := map[*string]interface{}{}
	for _, p := range schematic.parts {
		//log.Printf("Checking for the neighbors of the part at (%d,%d)", p[0], p[1])
		neighbors := schematic.Neighbors(p[0], p[1])

		for _, n := range neighbors {
			if n.number != nil {
				log.Printf("Got a numerical neighbor for (%d,%d): %s\t%p", p[0], p[1], *n.number, n.number)
				partNeighbors[n.number] = true
			}
		}
	}

	html := schematic.HTML()
	f, err := os.OpenFile("/tmp/schematic.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		log.Printf("Failed to write HTML visualization: %s", err)
	}
	_, _ = f.WriteString(html)

	log.Printf("Open HTML report in browser: file:///tmp/schematic.html")

	var sum int
	count := len(partNeighbors)
	log.Printf("Got %d labeled part numbers", count)
	for pointer := range partNeighbors {
		i, _ := strconv.ParseInt(*pointer, 10, 64)
		sum = sum + int(i)
	}
	return sum
}

func main() {
	log.Printf("Sum of schematic part numbers: %d", p1())
}
