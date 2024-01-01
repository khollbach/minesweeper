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

		row, col, err := parseLine(b, l)
		if err != nil {
			fmt.Println(err)
			continue
		}

		b.reveal(row, col)
	}
}

func parseLine(b board, line string) (row, col int, err error) {
	nrows, ncols := b.bounds()
	dims := []int{nrows, ncols}

	words := strings.Fields(line)
	if len(words) != 2 {
		return 0, 0, fmt.Errorf("Enter two numbers; e.g.: 1 2")
	}

	coords := make([]int, 2)
	for i, w := range words {
		n, err := strconv.Atoi(w)
		if err != nil {
			return 0, 0, fmt.Errorf("Expected a number, got %v", w)
		}
		if n < 0 {
			return 0, 0, fmt.Errorf("Negative number: %v", n)
		}
		if n >= dims[i] {
			rowCol := "row"
			if i == 1 {
				rowCol = "col"
			}
			return 0, 0, fmt.Errorf("Out of bounds. Largest %v coord is %v", rowCol, dims[i]-1)
		}
		coords[i] = n
	}

	return coords[0], coords[1], nil
}
