package main

import (
	"./tools"
	"github.com/hajimehoshi/ebiten"
	"log"
)

const (
	windowWidth  = 640
	windowHeight = 480
	scale        = 1
)

var (
	game         *tools.Game
	lastInputKey string // to prevent chain breaking
)

func update(surface *ebiten.Image) error {

	game.Surface = surface

	game.Bindings()

	//frame skip
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	switch game.CurrentScreen {
	case tools.TitleScreen:
		game.DrawTitleScreen()
	case tools.GameOverScreen:
		game.DrawGameOver()
	case tools.InGameScreen:
		game.DrawChain()
		game.DrawInGame(&lastInputKey)
	case tools.OptionsScreen:
	case tools.HallOfFameScreen:
	case tools.EnterInitialScreen:
	}

	game.DrawScore()
	game.DrawLives()
	game.DrawFPS()

	return nil
}

func main() {
	initGame()
	if err := ebiten.Run(update, game.Width, game.Height, scale, "Super Go Typing"); err != nil {
		log.Fatal(err)
	}
}

func initGame() {
	// create game
	game = tools.NewGame()
	game.Width = windowWidth
	game.Height = windowHeight

	game.CurrentScreen = tools.TitleScreen
}
