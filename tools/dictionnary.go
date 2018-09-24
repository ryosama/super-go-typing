package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

//Dictionnary contains words
type Dictionnary struct {
	Words []string
}

//NewDictionnary build a new dictionnary object from the path given (one word per line)
func NewDictionnary(path string) *Dictionnary {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	dic := new(Dictionnary)
	dic.Words = strings.Split(string(content), "\n")
	return dic
}

//GetRamdomWord extract a ramdom word in the dictionnary
func (dictionnary *Dictionnary) GetRamdomWord() string {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r1.Intn(len(dictionnary.Words))

	word := strings.ToUpper(strings.TrimSpace(dictionnary.Words[index]))
	fmt.Printf("Choose: len(%s):%d  index=%d\n", word, len(word), index)

	return word
}
