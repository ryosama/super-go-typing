package tools

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"os"
)

//Bindings checks keyboard interaction
func (game *Game) Bindings() {

	// if Alt+Enter, toogle fullscreen
	if ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if game.lastStatePressedKey["Alt+Enter"] == false {
			if ebiten.IsFullscreen() {
				ebiten.SetFullscreen(false)
			} else {
				ebiten.SetFullscreen(true)
			}
			game.lastStatePressedKey["Alt+Enter"] = true
		}
	} else {
		game.lastStatePressedKey["Alt+Enter"] = false
	}

	// escape --> quit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	// Entrer on Title screen --> start
	if game.CurrentScreen == TitleScreen && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		game.CurrentScreen = InGameScreen
		game.LaunchInGame()
	}

	// Entrer on Game over --> restart
	if game.CurrentScreen == GameOverScreen && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		game.CurrentScreen = TitleScreen
	}
}
