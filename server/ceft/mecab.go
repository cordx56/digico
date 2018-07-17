package ceft

import (
	"fmt"
	"strings"
	"github.com/shogo82148/go-mecab"
)

// Morp Morpheme struct (s: surface, f: feature)
type Morp struct {
	Surface string
	Feature []string
}
// ParsedSen ParsedSentence struct (plain: plain, morps: Morphemes)
type ParsedSen struct {
	Plain string
	Morps []Morp
}

// Parse Parse input string (str: input string)
func (p *ParsedSen) Parse(str string) int {
	p.Plain = str
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
			p.Morps = append(p.Morps, Morp{Surface: word[0], Feature: strings.Split(word[1], ",")})
		}
	}
	return 0
}

func MorpIsIndep(wclass string) bool {
	if (wclass == "助詞" || wclass == "助動詞" || wclass == "記号") {
		return false
	} else {
		return true
	}
}
func (m *Morp) IsIndep() bool {
	return MorpIsIndep(m.Feature[0])
}

func (p *ParsedSen) CountIndepMorps() int {
	mcount := 0
	for _, ele := range p.Morps {
		if ele.IsIndep() {
			mcount++
		}
	}
	return mcount
}
func (p *ParsedSen) CalcTF(m Morp) float32 {
	mcount := 0
	for _, ele := range p.Morps {
		if m.Surface == ele.Surface {
			mcount++
		}
	}
	return float32(mcount) / float32(p.CountIndepMorps())
}

func SplitText(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return (r == '。' || r == '．' || r == '\n' || r == '！' || r == '？')
	})
}
