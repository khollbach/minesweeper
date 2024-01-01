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
