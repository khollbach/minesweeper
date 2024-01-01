package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func main() {
	ebiten.SetWindowTitle("Minesweeper")

	b := newBoard()
	x, y := b.pixelDims()
	ebiten.SetWindowSize(x, y)

	if err := ebiten.RunGame(newGame(b)); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	board    board
	gameOver bool
}

func newGame(b board) *Game {
	return &Game{b, false}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.board.pixelDims()
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		i, j := pixelToTile(x, y)
		if g.board.inBounds(i, j) {
			g.board.reveal(i, j)
		}
	}

	state := g.board.gameOver()
	if state != inProgress {
		g.gameOver = true

		for _, row := range g.board {
			for j := range row {
				row[j].vis = revealed
			}
		}

		switch state {
		case win:
			fmt.Println("you won")
		case lose:
			fmt.Println("you lost")
		default:
			panic("unreachable")
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, row := range g.board {
		for j, tile := range row {
			var color color.RGBA
			var text string

			switch tile.vis {
			case noFlag:
				color = grey()
			case flag:
				color = red()
			case revealed:
				if tile.has_mine {
					color = black()
				} else {
					color = white()
					if n := g.board.neighboringMines(i, j); n != 0 {
						text = fmt.Sprint(n)
					}
				}
			default:
				panic("unreachable")
			}

			x, y := topLeft(i, j)
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(tileSize), float32(tileSize), color, false)
			ebitenutil.DebugPrintAt(screen, text, x, y)
		}
	}
}

// The side-length of a tile, in pixels.
const tileSize = 64

func (b board) pixelDims() (x, y int) {
	h, w := b.bounds()
	return w * tileSize, h * tileSize // Note the swap.
}

// Tile coords (row, col) -> pixel coords (x, y).
func topLeft(i, j int) (x, y int) {
	y, x = i, j // Note the swap.
	x *= tileSize
	y *= tileSize
	return x, y
}

// Pixel coords (x, y) -> containing tile (i, j).
func pixelToTile(x, y int) (i, j int) {
	x /= tileSize
	y /= tileSize
	return y, x // Note the swap.
}

func black() color.RGBA {
	return color.RGBA{0, 0, 0, 255}
}

func grey() color.RGBA {
	return color.RGBA{128, 128, 128, 255}
}

func white() color.RGBA {
	return color.RGBA{255, 255, 255, 255}
}

func red() color.RGBA {
	return color.RGBA{255, 0, 0, 255}
}
