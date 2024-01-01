package main

import (
	"fmt"
)

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

func (b board) bounds() (i, j int) {
	return len(b), len(b[0])
}

func (b board) inBounds(i, j int) bool {
	h, w := b.bounds()
	return 0 <= i && i < h && 0 <= j && j < w
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

// No trailing newline.
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
