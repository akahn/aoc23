package main

import (
	"bufio"
	"bytes"
	"embed"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input embed.FS

const (
	Red int = iota
	Green
	Blue
)

type BagPull struct {
	Red   int
	Green int
	Blue  int
}

type Game struct {
	ID    int
	Pulls []BagPull
}

func (g *Game) PossibleWith(red, green, blue int) bool {
	for _, p := range g.Pulls {
		if p.Red > red {
			return false
		}
		if p.Green > green {
			return false
		}
		if p.Blue > blue {
			return false
		}
	}

	return true
}

func (g *Game) Minimums() (minRed, minGreen, minBlue int) {
	for _, p := range g.Pulls {
		if p.Red > minRed {
			minRed = p.Red
		}
		if p.Green > minGreen {
			minGreen = p.Green
		}
		if p.Blue > minBlue {
			minBlue = p.Blue
		}
	}

	return
}

func parseGames() []Game {
	b, err := input.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("failed to open input: %s", err)
	}

	reader := bytes.NewReader(b)
	s := bufio.NewScanner(reader)

	games := []Game{}
	for {
		ok := s.Scan()
		if !ok {
			log.Printf("Got to end of file, breaking")
			break
		}

		line := s.Bytes()
		log.Printf("Got line %s", line)
		parts := strings.Split(string(line), ": ")
		intro := parts[0]
		game := parts[1]

		// Get the ID out
		introParts := strings.Split(intro, " ")
		id := introParts[1]

		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatalf("!")
		}
		g := Game{ID: int(i), Pulls: []BagPull{}}

		// Get the draws out
		draws := strings.Split(game, "; ")
		for _, draw := range draws {
			pull := BagPull{}
			colors := strings.Split(draw, ", ")
			for _, c := range colors {
				// split on space
				color := strings.Split(c, " ")
				switch color[1] {
				case "red":
					pull.Red = parseInt(color[0])
				case "blue":
					pull.Blue = parseInt(color[0])
				case "green":
					pull.Green = parseInt(color[0])
				default:
					log.Fatalf("Encountered unknown color: `%s` in `%s`", color[1], c)
				}
			}

			g.Pulls = append(g.Pulls, pull)
		}

		games = append(games, g)
	}

	return games
}

func p1(games []Game) int {
	possible := []int{}
	for _, g := range games {
		if g.PossibleWith(12, 13, 14) {
			possible = append(possible, g.ID)
		}
	}

	var sum int
	for _, p := range possible {
		sum = sum + p
	}

	return sum
}

func p2(games []Game) int {
	var sum int
	for _, game := range games {
		r, g, b := game.Minimums()
		power := r * g * b
		sum = sum + power
	}

	return sum
}

func main() {
	games := parseGames()
	possibleGames := p1(games)
	log.Printf("Sum of possible game IDs: %d", possibleGames)
	sumOfGamePowers := p2(games)
	log.Printf("Sum of game powers: %d", sumOfGamePowers)
}

func parseInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("Parsing int failed: %s", err)
	}

	return int(i)
}
