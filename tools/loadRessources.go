package tools

import (
	"log"
	"io/ioutil"
	"github.com/hajimehoshi/ebiten"
	"github.com/ryosama/go-sprite"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/font"
	"github.com/golang/freetype/truetype"
)

func loadSprite(path string) *sprite.Sprite {
	s := sprite.NewSprite()
	s.AddAnimation("default",path,	1, 1, ebiten.FilterDefault)
	s.Start()
	return s
}

func loadSound(path string) *audio.Player {
	audioContext, err := audio.NewContext(48000)
	if err != nil { log.Fatal(err) }
	
	f, err := ebitenutil.OpenFile(path) 		// load file
	if err != nil { log.Fatal(err) }
	d, err := vorbis.Decode(audioContext, f)  	// decode vorbis file
	if err != nil { log.Fatal(err) }

	audioPlayer, err := audio.NewPlayer(audioContext, d)
	if err != nil { log.Fatal(err) }
	return audioPlayer
}

func loadFont(path string, fontSize float64) font.Face {
	fontFile, err := ioutil.ReadFile(path)
	if err != nil { log.Fatal(err) }
	tt, err := truetype.Parse(fontFile)
	if err != nil {	log.Fatal(err) }

	return truetype.NewFace(tt,
				&truetype.Options{
					Size:    fontSize,
					DPI:     72,
					Hinting: font.HintingFull,
			})
}
