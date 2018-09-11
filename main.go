package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"./tools"
)

const (
	WINDOW_WIDTH  	= 640
	WINDOW_HEIGHT 	= 480
	SCALE        	= 1
)

var (
	game 			*tools.Game
	lastInputKey 	string	// to prevent chain breaking
)


func update(surface *ebiten.Image) error {

	game.Surface = surface

	game.Bindings()

	//frame skip
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	switch game.CurrentScreen {
		case tools.TITLE_SCREEN 		: 	game.DrawTitleScreen()
		case tools.GAMEOVER_SCREEN 		: 	game.DrawGameOver()
		case tools.INGAME_SCREEN 		: 	game.DrawChain() ; game.DrawInGame(&lastInputKey)
		case tools.OPTIONS_SCREEN 		: 	
		case tools.HALL_OF_FAME_SCREEN 	: 
		case tools.ENTER_INITIAL_SCREEN : 
	}

	game.DrawScore()
	game.DrawLives()
	game.DrawFPS()
	
	return nil
}



func main() {
	initGame()
	if err := ebiten.Run(update, game.Width, game.Height, SCALE, "Super Go Typing"); err != nil { log.Fatal(err) }
}

func initGame() {
	// create game
	game 			= tools.NewGame()
	game.Width 		= WINDOW_WIDTH
	game.Height 	= WINDOW_HEIGHT

	game.CurrentScreen = tools.TITLE_SCREEN
}