package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	b := newBoard()
	lines := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println(b)

		switch b.gameOver() {
		case lose:
			fmt.Println("You lost :(")
			return
		case win:
			fmt.Println("Congratulations! You won!!")
			return
		}

		fmt.Print("> ")
		ok := lines.Scan()
		if !ok {
			if err := lines.Err(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}
		l := lines.Text()

		words := strings.Fields(l)
		if len(words) != 2 {
			fmt.Println("Enter two numbers; e.g.: 1 2")
			continue
		}
		coords := make([]int, 2)
		var err error
		coords[0], err = strconv.Atoi(words[0])
		if err != nil {
			fmt.Println("Expected a number.", err)
			continue
		}
		coords[1], err = strconv.Atoi(words[1])
		if err != nil {
			fmt.Println("Expected a number.", err)
			continue
		}
		if coords[0] < 0 {
			fmt.Println("Negative number:", coords[0])
			continue
		}
		if coords[1] < 0 {
			fmt.Println("Negative number:", coords[1])
			continue
		}
		if coords[0] >= len(b) {
			fmt.Println("Out of bounds. Largest row coord is", len(b)-1)
			continue
		}
		if coords[1] >= len(b[0]) {
			fmt.Println("Out of bounds. Largest column coord is", len(b[0])-1)
			continue
		}

		b.reveal(coords[0], coords[1])
	}
}

func (b board) gameOver() gameOver {
	numRevealed := 0
	numMines := 0
	for _, row := range b {
		for _, tile := range row {
			if tile.vis == revealed {
				numRevealed++
			}
			if tile.has_mine {
				numMines++
			}
			if tile.vis == revealed && tile.has_mine {
				return lose
			}
		}
	}

	if numRevealed+numMines == len(b)*len(b[0]) {
		return win
	}
	return inProgress
}

type gameOver = int

const (
	inProgress gameOver = iota
	lose
	win
)

// Non-empty rectangle.
type board [][]tile

type tile struct {
	has_mine bool
	vis      visibility
}

type visibility = int

const (
	noFlag visibility = iota
	flag
	revealed
)

// If there's a flag, do nothing.
func (b board) reveal(i, j int) {
	if b[i][j].vis == flag {
		return
	}
	b[i][j].vis = revealed
}

// If the square is revealed, do nothing.
func (b board) toggleFlag(i, j int) {
	switch b[i][j].vis {
	case revealed:
	case flag:
		b[i][j].vis = noFlag
	case noFlag:
		b[i][j].vis = flag
	default:
		panic("unreachable")
	}
}

func newBoard() board {
	o := tile{false, noFlag}
	x := tile{true, noFlag}

	return [][]tile{
		{o, o, o, o, x, x, o, o, o, o},
		{x, o, o, o, o, x, o, o, x, o},
		{o, o, o, o, o, o, o, o, o, o},
		{x, o, o, o, o, o, o, o, o, o},
		{o, o, o, o, o, o, o, o, o, o},
		{o, o, o, o, o, o, x, x, o, o},
		{o, o, o, o, x, o, o, o, o, o},
		{o, o, o, o, o, o, o, o, x, o},
	}
}

func (b board) String() string {
	s := ""
	for i := range b {
		if i > 0 {
			s += "\n"
		}
		for j := range b[i] {
			switch b[i][j].vis {
			case noFlag:
				s += "."
			case flag:
				s += "F"
			case revealed:
				if b[i][j].has_mine {
					s += "*"
				} else {
					n := b.neighboringMines(i, j)
					if n == 0 {
						s += " "
					} else {
						s += fmt.Sprint(n)
					}
				}
			default:
				panic("unreachable")
			}
		}
	}
	return s
}

func (b board) neighboringMines(i, j int) int {
	out := 0
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			orig := di == 0 && dj == 0
			i2 := i + di
			j2 := j + dj
			if !orig && b.inBounds(i2, j2) && b[i2][j2].has_mine {
				out++
			}
		}
	}
	return out
}

func (b board) bounds() (i, j int) {
	return len(b), len(b[0])
}

func (b board) inBounds(i, j int) bool {
	h, w := b.bounds()
	return 0 <= i && i < h && 0 <= j && j < w
}
