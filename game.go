package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ------------------------------- DATA -------------------------------

type Game struct {
	ImgSize               int
	GridWidth, GridHeight int
	NumMines              int
	RevealedTiles         int
	Sprites               []*ebiten.Image
	SpritesGrid           [][]int
	LogicGrid             [][]int
	MinesPos              []Pos
	FlagsPos              map[Pos]any

	PrevMouseLeftPressed  bool
	PrevMouseRightPressed bool
}

type Pos struct {
	i, j int
}

// ------------------------------- INIT / RESET -------------------------------

func Init(gridWidth, gridHeight, imgSize, numMines int) *Game {
	g := &Game{}

	// Image Size
	g.ImgSize = imgSize

	// Grid Size
	g.GridWidth, g.GridHeight = gridWidth, gridHeight

	// Number of mines
	g.NumMines = numMines

	g.RevealedTiles = 0

	// Sprites
	g.Sprites = make([]*ebiten.Image, NUM_SPRITES)
	for i, spritePath := range SpritesPaths {
		g.Sprites[i] = loadImage(spritePath)
	}

	// Sprites Grid
	g.SpritesGrid = make([][]int, gridHeight)
	for i := range len(g.SpritesGrid) {
		g.SpritesGrid[i] = make([]int, gridWidth)
		for j := range len(g.SpritesGrid[i]) {
			g.SetSprite(i, j, HIDDEN)
		}
	}

	// Logic Grid
	g.LogicGrid = make([][]int, gridHeight)
	for i := range g.LogicGrid {
		g.LogicGrid[i] = make([]int, gridWidth)
	}
	g.GenerateMines()
	g.GenerateNumbers()

	g.FlagsPos = make(map[Pos]any)

	g.PrevMouseLeftPressed = false
	g.PrevMouseRightPressed = false

	return g
}

func (g *Game) Reset() {

	// Sprites Grid
	g.SpritesGrid = make([][]int, g.GridHeight)
	for i := range len(g.SpritesGrid) {
		g.SpritesGrid[i] = make([]int, g.GridWidth)
		for j := range len(g.SpritesGrid[i]) {
			g.SetSprite(i, j, HIDDEN)
		}
	}

	// Logic Grid
	g.LogicGrid = make([][]int, g.GridHeight)
	for i := range g.LogicGrid {
		g.LogicGrid[i] = make([]int, g.GridWidth)
	}
	g.GenerateMines()
	g.GenerateNumbers()

	g.RevealedTiles = 0

	g.FlagsPos = make(map[Pos]any)

	gameRunning = true
	gameWon = false
	gameLost = false
}

func (g *Game) GenerateMines() {
	g.MinesPos = make([]Pos, 0, g.NumMines)
	total := g.GridWidth * g.GridHeight
	tiles := make([]int, total)
	for i := 0; i < total; i++ {
		tiles[i] = i
	}
	rand.Shuffle(total, func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})

	for idx, val := range tiles {
		if idx == g.NumMines {
			break
		}
		i, j := val/g.GridWidth, val%g.GridWidth
		g.SetLogic(i, j, MINE)
		g.MinesPos = append(g.MinesPos, Pos{i, j})
	}
}

func (g *Game) GenerateNumbers() {
	for i := range g.GridHeight {
		for j := range g.GridWidth {
			if g.GetLogic(i, j) != MINE {
				g.SetLogic(i, j, len(g.GetValidNeighborsF(i, j, func(p Pos) bool {
					return g.GetLogic(p.i, p.j) == MINE
				})))
			}
		}
	}
}

// ------------------------------- GETTERS / SETTERS -------------------------------

func (g *Game) SetSprite(i, j, val int) {
	g.SpritesGrid[i][j] = val
}

func (g *Game) SetLogic(i, j, val int) {
	g.LogicGrid[i][j] = val
}

func (g *Game) GetSprite(i, j int) int {
	return g.SpritesGrid[i][j]
}

func (g *Game) GetLogic(i, j int) int {
	return g.LogicGrid[i][j]
}

// ------------------------------- EBITEN LOGIC -------------------------------

func (g *Game) Update() error {
	i, j := g.GetMouseClickGridPos()
	if i < 0 || j < 0 || i >= g.GridHeight || j >= g.GridWidth {
		return nil
	}

	currSprite := g.GetSprite(i, j)

	// Always called even if not used to handle mouse release logic
	_ = g.MouseLeftPressed()
	mouseLeftReleased := g.MouseLeftReleased()
	_ = g.MouseRightPressed()
	mouseRightReleased := g.MouseRightReleased()

	if gameRunning {

		if mouseLeftReleased {
			if g.GetSprite(i, j) == HIDDEN {
				g.Reveal(i, j)

				if g.GetLogic(i, j) == MINE {
					g.HandleLoose(i, j)
				} else if g.GetLogic(i, j) == EMPTY {
					g.ExpandEmpty(i, j)
				}
			}
			if ONE <= g.GetSprite(i, j) && g.GetSprite(i, j) <= EIGHT {
				g.RevealWithFlags(i, j)
			}
		}

		if mouseRightReleased {
			if currSprite == HIDDEN {
				g.SetSprite(i, j, FLAG)
				g.FlagsPos[Pos{i, j}] = nil
			} else if currSprite == FLAG {
				g.SetSprite(i, j, HIDDEN)
				delete(g.FlagsPos, Pos{i, j})
			}
		}

	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Reset()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for i := range len(g.SpritesGrid) {
		for j := range len(g.SpritesGrid[i]) {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(float64(g.ImgSize)/SPRITE_SIZE, float64(g.ImgSize)/SPRITE_SIZE)
			opts.GeoM.Translate(float64(j*g.ImgSize), float64(i*g.ImgSize))
			screen.DrawImage(g.Sprites[g.GetSprite(i, j)], opts)
		}
	}

	if !gameRunning {
		if gameWon {
			ebitenutil.DebugPrint(screen, "You won!, press SPACE to restart")
		}
		if gameLost {
			ebitenutil.DebugPrint(screen, "Game over, press SPACE to restart")
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.GridWidth * g.ImgSize, g.GridHeight * g.ImgSize
}

// ------------------------------- MOUSE LOGIC -------------------------------

func (g *Game) GetMouseClickGridPos() (int, int) {
	x, y := ebiten.CursorPosition()
	return y / g.ImgSize, x / g.ImgSize
}

func (g *Game) MouseLeftPressed() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.PrevMouseLeftPressed = true
		return true
	}
	return false
}

func (g *Game) MouseLeftReleased() bool {
	if g.PrevMouseLeftPressed && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.PrevMouseLeftPressed = false
		return true
	}
	return false
}

func (g *Game) MouseRightPressed() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.PrevMouseRightPressed = true
		return true
	}
	return false
}

func (g *Game) MouseRightReleased() bool {
	if g.PrevMouseRightPressed && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.PrevMouseRightPressed = false
		return true
	}
	return false
}

// ------------------------------- GAME LOGIC -------------------------------

func (g *Game) HandleWin() {
	toReveal := g.GridWidth*g.GridHeight - g.RevealedTiles
	if toReveal != g.NumMines {
		return
	}
	for _, p := range g.MinesPos {
		if g.GetSprite(p.i, p.j) != HIDDEN && g.GetSprite(p.i, p.j) != FLAG {
			return
		}
	}

	for _, p := range g.MinesPos {
		g.SetSprite(p.i, p.j, FLAG)
	}

	gameRunning = false
	gameWon = true
}

func (g *Game) HandleLoose(i, j int) {

	for _, p := range g.MinesPos {
		if p.i == i && p.j == j {
			g.SetSprite(i, j, MINE_EXPLODE)
			continue
		}

		if g.GetSprite(p.i, p.j) != FLAG {
			g.SetSprite(p.i, p.j, MINE)
		}
	}

	for p := range g.FlagsPos {
		if g.GetLogic(p.i, p.j) != MINE {
			g.SetSprite(p.i, p.j, FLAG_WRONG)
		}
	}

	gameRunning = false
	gameLost = true
}

func (g *Game) ExpandEmpty(i, j int) {
	noMineNeighbors := g.BFS(i, j)
	for _, nb := range noMineNeighbors {
		if g.GetSprite(nb.i, nb.j) == HIDDEN {
			g.Reveal(nb.i, nb.j)
		}
	}
}

func (g *Game) RevealWithFlags(i, j int) {
	neighborsFlags := g.GetValidNeighborsF(i, j, func(p Pos) bool {
		return g.GetSprite(p.i, p.j) == FLAG
	})

	if len(neighborsFlags) != g.GetLogic(i, j) {
		return
	}

	revealedNeighbors := g.GetValidNeighborsF(i, j, func(p Pos) bool {
		return g.GetSprite(p.i, p.j) == HIDDEN
	})

	for _, nb := range revealedNeighbors {
		k, l := nb.i, nb.j
		g.Reveal(k, l)

		if g.GetLogic(k, l) == EMPTY {
			g.ExpandEmpty(i, j)
		}

		if g.GetLogic(k, l) == MINE {
			g.HandleLoose(k, l)
		}
	}
}

// ------------------------------- UTILITY -------------------------------

func (g *Game) Reveal(i, j int) {
	if g.GetSprite(i, j) == HIDDEN {
		g.SetSprite(i, j, g.GetLogic(i, j))
		g.RevealedTiles++
		g.HandleWin()
	}
}

func (g *Game) GetValidNeighborsF(i, j int, f func(Pos) bool) []Pos {
	dirs := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	nbs := make([]Pos, 0, 8)
	for _, d := range dirs {
		nb := Pos{i + d[0], j + d[1]}
		if !g.IsValid(nb) {
			continue
		}
		if f(nb) {
			nbs = append(nbs, nb)
		}
	}
	return nbs
}

func (g *Game) IsValid(p Pos) bool {
	return 0 <= p.i && p.i < g.GridHeight && 0 <= p.j && p.j < g.GridWidth
}

func (g *Game) BFS(i, j int) []Pos {
	initPos := Pos{i, j}

	// Queue
	queue := make([]Pos, 0, 32)

	// Visited
	visited := make(map[Pos]bool)

	// Border
	border := make(map[Pos]bool)

	// Enqueue + Visit ij
	queue = append(queue, initPos)
	visited[initPos] = true

	// While Queue
	for len(queue) > 0 {

		currPos := queue[0]
		queue = queue[1:]

		for _, nb := range g.GetValidNeighborsF(currPos.i, currPos.j, func(p Pos) bool {
			return g.GetLogic(p.i, p.j) != MINE
		}) {
			if _, ok := visited[nb]; !ok {
				k, l := nb.i, nb.j
				if g.GetLogic(k, l) == EMPTY {
					queue = append(queue, nb)
					visited[nb] = true
				} else {
					if _, ok := border[nb]; !ok {
						border[nb] = true
					}
				}
			}
		}
	}

	// Return Visited + Border
	ret := make([]Pos, len(visited)-1+len(border))
	idx := 0
	for k := range visited {
		if k != initPos {
			ret[idx] = k
			idx++
		}
	}
	for k := range border {
		ret[idx] = k
		idx++
	}

	return ret
}
