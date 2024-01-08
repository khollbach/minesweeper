package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func main() {
	ebiten.SetWindowTitle("Minesweeper")

	b := newBoard()

	// Scale up the window to make it easier to see.
	x, y := b.pixelDims()
	x *= 4
	y *= 4
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

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		i, j := pixelToTile(x, y)
		if g.board.inBounds(i, j) {
			g.board.toggleFlag(i, j)
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
			var sprite string

			switch tile.vis {
			case noFlag:
				sprite = "blank"
			case flag:
				sprite = "flag"
			case revealed:
				if tile.has_mine {
					sprite = "mine"
				} else {
					sprite = "blank-pressed"
					if n := g.board.neighboringMines(i, j); n != 0 {
						sprite = fmt.Sprint(n)
					}
				}
			default:
				panic("unreachable")
			}

			// TODO: it's silly to read the files from disk every frame;
			// we should do this once and save them in memory.
			reader, err := os.Open("sprites/" + sprite + ".png")
			if err != nil {
				log.Fatalf("couldn't open file for sprite %q %v", sprite, err)
			}
			image, fmt, err := image.Decode(reader)
			if err != nil {
				log.Fatalf("error decoding sprite %q %q %v", sprite, fmt, err)
			}
			if fmt != "png" {
				log.Fatalf("expected png got %q", fmt)
			}

			x, y := topLeft(i, j)
			im := ebiten.NewImageFromImage(image)
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(im, &op)
		}
	}
}

// The side-length of a tile, in pixels.
const tileSize = 16

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
