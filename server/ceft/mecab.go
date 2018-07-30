package ceft

import (
	"fmt"
	"strings"
	"github.com/shogo82148/go-mecab"
)


// ParseMeCab Parse input string using MeCab (str: input string)
func (p *Sentence) ParseMeCab(str string) int {
	p.Plain = str
	p.MorphemeMap = make(map[string]*Morpheme)
	tagger, err := mecab.New(map[string]string{"dicdir": "/usr/lib/mecab/dic/mecab-ipadic-neologd"})
	defer tagger.Destroy()
	if (err != nil) {
		fmt.Println(err)
		return 1
	}
	result, err := tagger.Parse(str)
	if (err != nil) {
		fmt.Println(err)
		return 1
	}
	for _, ele := range strings.Split(result, "\n") {
		word:= strings.Split(ele, "\t")
		if (len(word) > 1) {
			p.Morphemes = append(p.Morphemes, Morpheme{Surface: word[0], Feature: strings.Split(word[1], ",")})
			p.MorphemeMap[word[0]] = &p.Morphemes[len(p.Morphemes) - 1]
		}
	}
	return 0
}

func MorphemeIsIndep(wclass string) bool {
	if (wclass == "助詞" || wclass == "助動詞" || wclass == "記号") {
		return false
	}
	return true
}
func (m *Morpheme) IsIndep() bool {
	return MorphemeIsIndep(m.Feature[0])
}

func (p *Sentence) CountIndepMorphemes() int {
	mcount := 0
	for _, ele := range p.Morphemes {
		if ele.IsIndep() {
			mcount++
		}
	}
	return mcount
}
