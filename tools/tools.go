package tools

import (
	"image/color"
	"strings"
	"reflect"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"golang.org/x/image/font"
)


//////////////////////////////////////////////////////////////////////////////////
///////////////////////////// Static METHODS /////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////



func DrawRectOutFill(surface *ebiten.Image, x, y , w, h float64, c color.Color) {
	x1 := x+w
	y1 := y+h
	ebitenutil.DrawLine(surface, x,  y, x1,  y, c) 		// top
	ebitenutil.DrawLine(surface, x, y1, x1, y1, c)		// bottom
	ebitenutil.DrawLine(surface, x,  y,  x, y1, c)		// left
	ebitenutil.DrawLine(surface,x1, y,  x1, y1, c)		// right
}


func GetSizeOfText(s string, f font.Face) (float64,float64) {
	rect,_ := font.BoundString(f,s)
	width  := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	return (float64(width) / (1 << 6))+1, (float64(height) / (1 << 6))+1
}

func GetInputChar() string {
	keys := ebiten.InputChars()
	activeLetter := ""
	if len(keys)>0 { // if key is pressed
		activeLetter = strings.ToUpper(string(keys[0]))
	}
	return activeLetter
}

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

