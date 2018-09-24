package tools

import (
	"image/color"
	"reflect"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"golang.org/x/image/font"
)

//DrawRectOutFill draw a rectangle outfill at x,y coordonnates
func DrawRectOutFill(surface *ebiten.Image, x, y, w, h float64, c color.Color) {
	x1 := x + w
	y1 := y + h
	ebitenutil.DrawLine(surface, x, y, x1, y, c)   // top
	ebitenutil.DrawLine(surface, x, y1, x1, y1, c) // bottom
	ebitenutil.DrawLine(surface, x, y, x, y1, c)   // left
	ebitenutil.DrawLine(surface, x1, y, x1, y1, c) // right
}

//GetSizeOfText return the size of a text un pixel with the given font face
func GetSizeOfText(s string, f font.Face) (float64, float64) {
	rect, _ := font.BoundString(f, s)
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	return (float64(width) / (1 << 6)) + 1, (float64(height) / (1 << 6)) + 1
}

//GetInputChar return the current caracters input on keyboard
func GetInputChar() string {
	keys := ebiten.InputChars()
	activeLetter := ""
	if len(keys) > 0 { // if key is pressed
		activeLetter = strings.ToUpper(string(keys[0]))
	}
	return activeLetter
}

//InArray return true and the index of a value in a given array (false and -1 otherwise)
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}
