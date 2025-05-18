package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// type Tile struct {
// 	i, j   int
// 	sprite int
// }

const EMPTY = 0
const ONE = 1
const TWO = 2
const THREE = 3
const FOUR = 4
const FIVE = 5
const SIX = 6
const SEVEN = 7
const EIGHT = 8
const HIDDEN = 9
const MINE = 10
const MINE_EXPLODE = 11
const FLAG = 12
const FLAG_WRONG = 13

const NUM_SPRITES = 14
const SPRITE_SIZE = 128

var SpritesPaths = []string{
	"assets/empty.png",
	"assets/1.png",
	"assets/2.png",
	"assets/3.png",
	"assets/4.png",
	"assets/5.png",
	"assets/6.png",
	"assets/7.png",
	"assets/8.png",
	"assets/hidden.png",
	"assets/mine.png",
	"assets/mine_explode.png",
	"assets/flag.png",
	"assets/flag_wrong.png",
}

var gameRunning = true
var gameWon = false
var gameLost = false

//var clickedLeftPos *Tile = nil

func loadImage(path string) *ebiten.Image {

	// Open file
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error: could not open file %s: %v", path, err)
	}
	defer f.Close()

	// Decode Image
	imageData, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("Error when decoding image: %v", err)
	}

	return ebiten.NewImageFromImage(imageData)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func Debug(s string) {
	fmt.Println(s)
}
