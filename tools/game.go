package tools 

//////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// Game METHODS /////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/audio"
	
	"golang.org/x/image/font"

	"github.com/ryosama/go-sprite"
)

const (
	INGAME_SCREEN = iota
	GAMEOVER_SCREEN
	OPTIONS_SCREEN
	TITLE_SCREEN
	HALL_OF_FAME_SCREEN
	ENTER_INITIAL_SCREEN

	INITIAL_SPEED 		= 0.2
	INCREMEMENT_SPEED 	= 0.1
	LIVES				= 5
)

type Game struct {
	Dictionnary 			*Dictionnary
	Width, Height 			int
	Words					[3]*Word
	Score,Chain				int
	Speed					float64
	Surface 				*ebiten.Image
	Lives, MaxLives			int
	CurrentScreen			int
	lastStatePressedKey		map[string]bool
	sprites 				map[string]*sprite.Sprite
	audioPlayers 			map[string]*audio.Player
	fonts 					map[string]font.Face
	colors 					map[string]color.RGBA
}


func NewGame() *Game {
	game := new(Game)
	
	game.lastStatePressedKey 	= make(map[string]bool)
	game.sprites 				= make(map[string]*sprite.Sprite)
	game.audioPlayers 			= make(map[string]*audio.Player)
	game.fonts 					= make(map[string]font.Face)
	game.colors					= make(map[string]color.RGBA)	

	// gfx
	game.sprites["heart-full"]  	= loadSprite("gfx/heart_full.png")
	game.sprites["heart-empty"] 	= loadSprite("gfx/heart_empty.png")
	game.sprites["title-screen"]  	= loadSprite("gfx/title_screen.png")

	// sfx
	game.audioPlayers["explosion"] 	= loadSound("sfx/misc-explosion.ogg")

	// fonts
	game.fonts["ui"] 		= loadFont("fonts/99 3Dventure.ttf"	, 30)
	game.fonts["ui-big"] 	= loadFont("fonts/99 3Dventure.ttf"	, 70)
	game.fonts["word"] 		= loadFont("fonts/00 impact.ttf"	, 20)

	// load dictionnaries
	//game.Dictionnary= NewDictionnary("dictionnaries/accent.txt")
	//game.Dictionnary= NewDictionnary("dictionnaries/a.txt")
	game.Dictionnary = NewDictionnary("dictionnaries/french2.txt")

	// load colors
	game.colors["white"]			= color.RGBA{ R:255, G:255, B:255, A:255 }
	game.colors["black"]			= color.RGBA{ R:0  , G:0  , B:0  , A:255 }
	game.colors["red"]				= color.RGBA{ R:255, G:0  , B:0  , A:255 }
	game.colors["green"]			= color.RGBA{ R:0  , G:255, B:0  , A:255 }
	game.colors["blue"]				= color.RGBA{ R:0  , G:0  , B:255, A:255 }
	game.colors["violet"]			= color.RGBA{ R:255, G:0  , B:255, A:255 }
	game.colors["yellow"]			= color.RGBA{ R:255, G:255, B:0  , A:255 }
	game.colors["turquoise"]		= color.RGBA{ R:0  , G:255, B:255, A:255 }
	game.colors["grey"]				= color.RGBA{ R:128, G:128, B:128, A:255 }

	game.colors["ui"] 				= color.RGBA{ R:105, G:82, 	B:186, A:255 }
	game.colors["word"]				= game.colors["turquoise"]
	game.colors["active-letter"]	= game.colors["red"]

	return game
}


func (game *Game) LaunchInGame() {
	game.Words = [3]*Word{}
	game.Lives 					= LIVES
	game.MaxLives 				= LIVES
	game.Speed 					= INITIAL_SPEED
	game.Score = 0
	game.Chain = 0

	game.NewWord(0)

	timer1 := time.NewTimer(2000 * time.Millisecond)
	go func() {
		<-timer1.C
		game.NewWord(1)
    }()

    timer2 := time.NewTimer(4000 * time.Millisecond)
	go func() {
		<-timer2.C
		game.NewWord(2)
    }()
}


func (game *Game) Draw(word *Word) {
	surface := game.Surface

	if !word.Fail {
		if word.Success { // draw +1 moving up
			text.Draw(surface, fmt.Sprintf("+%d",word.WordScore), game.fonts["ui"], int(word.X), int(word.Y), game.colors["ui"])
			word.Y -= 0.2 // move upward

		} else {	// draw the word
			x := word.X
			startOfWord := word.Text[0 : word.CurrentLetterIndex]
			w,_ := GetSizeOfText(startOfWord, game.fonts["word"])
			text.Draw(surface, startOfWord, game.fonts["word"], int(word.X), int(word.Y), game.colors["word"])
			
			x += w
			nextLetterToTape,numberOfOctet := word.GetLetterAtIndex(word.CurrentLetterIndex)
			text.Draw(surface, nextLetterToTape, game.fonts["word"], int(x), int(word.Y), game.colors["active-letter"])
			w,_ = GetSizeOfText(nextLetterToTape, game.fonts["word"])

			x += w
			endOfWord := word.Text[word.CurrentLetterIndex+numberOfOctet : ] // to the end
			text.Draw(surface, endOfWord, game.fonts["word"], int(x), int(word.Y), game.colors["word"])
		}

	} else { // draw fail animation
		word.Explosion.Draw(surface)
	}
}


func (game *Game) DrawInGame(lastInputKey *string) {

	activeLetter := GetInputChar()

	activeWord := false
	if game.GetActiveWordIndex() > -1 {
		activeWord = true
	}

	nextValidLetters := game.GetNextValidLetters()

	for index,word := range game.Words {
		if word != nil {

			// too late --> hide word and draw explosion
			if word.X + float64(word.Width) > float64(game.Width) && !word.Fail {
				game.CreateFail(index) // update word.Fail to false
				game.Chain = 0
				game.DestroyWord(index)
			}

			// currently active
	 		if word.IsActive && word.CurrentLetterIndex >= len(word.Text) && !word.Success { // success
	 			game.CreateSuccess(index) // update word.Success to true
	 			game.DestroyWord(index)
	 			game.Score += game.Chain
				game.Speed += INCREMEMENT_SPEED
	 		}

			// Test if activeLetter is the next letter to tape
			if len(activeLetter)>0 {
				if !word.IsDestroyed {
					if word.IsActive || !activeWord {
						nextLetterToTape,numberOfOctet := word.GetLetterAtIndex(word.CurrentLetterIndex)

						if inArray,_ := InArray(activeLetter, nextValidLetters) ; !inArray {
							game.Chain = 0
						}

						if activeLetter == nextLetterToTape {
							game.Chain++
							word.CurrentLetterIndex += numberOfOctet
							word.IsActive = true
							activeWord = true
						}
					}
				}
				*lastInputKey = activeLetter
			}

			game.Draw(word) // draw fail or success or normal animation

			// move word to the right
			if !word.Success {
				word.X += game.Speed
			}

		} // end word is not null
	} // end for each words
}

func (game *Game) NewWord(index int) {
	word := new(Word)

	// avoid word starting with the same letter
getValidWord:
	word.Text = game.Dictionnary.GetRamdomWord()
	for _,otherWords := range game.Words {
		if otherWords != nil {
			l1,_ := otherWords.GetLetterAtIndex(0)
			l2,_ := word.GetLetterAtIndex(0)
			if l1 == l2 { // same first letter
				goto getValidWord
			}
		}
	}

	word.X = 0
	word.Width, word.Height = GetSizeOfText(word.Text , game.fonts["word"])

	word.Y = float64(game.Height / 6 * (index*2 + 1))

	game.Words[index] = word
}

func (game *Game) DestroyWord(index int) {
	// destroy word after 1 second

	game.Words[index].IsActive 		= false // no more active
	game.Words[index].IsDestroyed 	= true  // no more active

	timer := time.NewTimer(1 * time.Second)
	go func() {
		<-timer.C
		game.Words[index] = nil
		game.NewWord(index)
	}()
}

func (game *Game) GetActiveWordIndex() int {
	for index,word := range game.Words {
		if word != nil {
			if word.IsActive {
				return index
			}
		}
	}
	return -1
}


func (game *Game) CreateFail(index int) {
	word := game.Words[index]
	word.Fail = true

	// choose an animation for explosion
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	explosion := sprite.NewSprite()
	explosionIndex := int(r1.Intn(3)+1) // bewteen 1 et 3
	switch (explosionIndex) {
		case 1:explosion.AddAnimation("explosion","gfx/explosion1.png", 500, 5, ebiten.FilterDefault)
		case 2:explosion.AddAnimation("explosion","gfx/explosion2.png", 500, 7, ebiten.FilterDefault)
		case 3:explosion.AddAnimation("explosion","gfx/explosion3.png", 500, 9, ebiten.FilterDefault)
	}

	explosion.CurrentAnimation = "explosion"
	explosion.X = word.X + (word.Width / 2 - explosion.GetWidth()/2) // in the middle of the word
	explosion.Y = word.Y - word.Height*2
	explosion.RunOnce(AfterRunOnce)

	word.Explosion = explosion

	// play the explosion sound
	game.audioPlayers["explosion"].Rewind()
	game.audioPlayers["explosion"].Play()

	game.Lives--

	if game.Lives <= 0 { // just lose
		game.CurrentScreen = GAMEOVER_SCREEN
	}
}

func (game *Game) CreateSuccess(index int) {
	game.Words[index].Success 	= true
	game.Words[index].WordScore = game.Chain
}

func (game *Game) DrawScore() {
	text.Draw(game.Surface, "Score: "+fmt.Sprintf("%06d",game.Score), game.fonts["ui"], game.Width - 250 , 20 , game.colors["ui"])
}

func (game *Game) DrawChain() {
	text.Draw(game.Surface, "Chain "+fmt.Sprintf("%v",game.Chain), game.fonts["ui"], game.Width - 200 , 200 , game.colors["ui"])
}


func (game *Game) DrawLives() {
	surface := game.Surface
	for i:=1 ; i<=game.MaxLives ; i++ {
		if i<=game.Lives {
			game.sprites["heart-full"].Position(float64(i)*20, float64(game.Height) - 20)
			game.sprites["heart-full"].Draw(surface)
		} else {
			game.sprites["heart-empty"].Position(float64(i)*20, float64(game.Height) - 20)
			game.sprites["heart-empty"].Draw(surface)
		}
	}
}

func (game *Game) DrawGameOver() {
	s := "GAME OVER"
	w,h := GetSizeOfText(s, game.fonts["ui-big"])
	text.Draw(game.Surface, s, game.fonts["ui-big"], game.Width/2 - int(w/2) , game.Height/2 - int(h/2) , game.colors["ui"])

	s = "Press 'Entrer' to restart"
	w,h = GetSizeOfText(s, game.fonts["ui"])
	text.Draw(game.Surface, s, game.fonts["ui"], game.Width/2 - int(w/2) , game.Height/2 + 30 - int(h/2) , game.colors["ui"])
}


func (game *Game) DrawTitleScreen() {
	game.sprites["title-screen"].Position(0,0)
	game.sprites["title-screen"].Draw(game.Surface)

	s := "Press 'Entrer' to start"
	w,h := GetSizeOfText(s, game.fonts["ui"])
	text.Draw(game.Surface, s, game.fonts["ui"], game.Width/2 - int(w/2) , game.Height/2 + 80 - int(h/2) , game.colors["ui"])
}


func (game *Game) GetNextValidLetters() []string {
	nextValidsLetters := []string{}

	if game.GetActiveWordIndex() > -1 { // only one valid letter
		word := game.Words[game.GetActiveWordIndex()]
		nextLetter,_ := word.GetLetterAtIndex(word.CurrentLetterIndex)
		nextValidsLetters = append( nextValidsLetters, nextLetter )

	} else {	// many valid letters
		for _,word := range game.Words {
			if word != nil {
				nextLetter,_ := word.GetLetterAtIndex(word.CurrentLetterIndex)
				nextValidsLetters = append( nextValidsLetters, nextLetter )
			}
		}
	}
	return nextValidsLetters
}


func (game *Game) DrawFPS() {
	// display FPS and other stuff
	ebitenutil.DebugPrint(game.Surface,fmt.Sprintf("FPS:%f",ebiten.CurrentFPS()))
}


func AfterRunOnce(s *sprite.Sprite) {
	//s = nil // destroy animation
}