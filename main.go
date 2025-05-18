package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	Debug("")

	gridWidth := flag.Int("w", 30, "grid width")
	gridHeight := flag.Int("h", 16, "grid height")
	numMines := flag.Int("m", 99, "number of mines")
	flag.Parse()

	if *numMines > *gridWidth**gridHeight {
		fmt.Println("That much mines!? It doesn't even fit on the grid!")
		os.Exit(1)
	}

	monitorWidth, monitorHeight := ebiten.Monitor().Size()
	imgSize := int(min(0.9*float64(monitorWidth)/float64(*gridWidth), 0.9*float64(monitorHeight)/float64(*gridHeight)))

	g := Init(*gridWidth, *gridHeight, imgSize, *numMines)

	ebiten.SetWindowSize(*gridWidth*imgSize, *gridHeight*imgSize)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowTitle("Minesweeper")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
