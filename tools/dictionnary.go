package tools

import (
	"log"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
	"fmt"
)

//////////////////////////////////////////////////////////////////////////////////
///////////////////////////// Dictionnary METHODS ////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////

type Dictionnary struct {
	Words 	[]string	
}

func NewDictionnary(path string) *Dictionnary {
	content, err := ioutil.ReadFile(path)
	if err != nil { log.Fatal(err) }

	dic := new(Dictionnary)
	dic.Words = strings.Split(string(content),"\n")
	return dic
}

func (dictionnary *Dictionnary) GetRamdomWord() string {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r1.Intn(len(dictionnary.Words))

	word := strings.ToUpper(strings.TrimSpace(dictionnary.Words[index]))
	fmt.Printf("Choose: len(%s):%d  index=%d\n",word,len(word),index)

	return word
}