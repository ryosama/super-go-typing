package tools

import (
	"github.com/ryosama/go-sprite"
)

//Word contains information about a word currently displayed
type Word struct {
	X, Y, Width, Height           float64
	Text                          string
	Explosion                     *sprite.Sprite
	Fail, Success                 bool
	CurrentLetterIndex, WordScore int
	IsActive, IsDestroyed         bool
}

//GetLetterAtIndex returns the letter at given index
func (word *Word) GetLetterAtIndex(index int) (string, int) {
	if index >= len(word.Text) {
		return "", 0
	}

	jumpForUTF8 := 1
	if word.Text[index] == 195 {
		jumpForUTF8++
	}

	return word.Text[index : index+jumpForUTF8], jumpForUTF8
}
