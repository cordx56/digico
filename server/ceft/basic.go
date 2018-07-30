package ceft

import (
	"math"
	"strings"
)


// Morpheme struct
type Morpheme struct {
	Surface string
	Feature []string
	Value float32
}
// Sentence parsed sentence struct
type Sentence struct {
	Plain string
	Morphemes []Morpheme
	MorphemeMap map[string]*Morpheme
}
// Document document struct
type Document struct {
	Plain string
	Sentences []Sentence
	MorphemeCount map[string]int
}
// Documents documents struct
type Documents struct {
	IDF map[string]float32
	Documents []Document
}


// CalcTF - calc TF of Morpheme in Sentence
func (p *Sentence) CalcTF(m Morpheme) float32 {
	mcount := 0
	for _, ele := range p.Morphemes {
		if m.Surface == ele.Surface {
			mcount++
		}
	}
	return float32(mcount) / float32(len(p.Morphemes))
}

// CalcIDF - calc IDF of all Morphemes in Documents
func (ds *Documents) CalcIDF() {
	ds.IDF = make(map[string]float32)
	for _, doc := range ds.Documents {
		morpCount := make(map[string]int)
		for _, sen := range doc.Sentences {
			for _, morp := range sen.Morphemes {
				morpCount[morp.Surface]++
			}
		}
		for morp, c := range morpCount {
			if (c > 0) { ds.IDF[morp]++ }
		}
	}
	for morp, v := range ds.IDF {
		ds.IDF[morp] = float32(math.Log2(float64(len(ds.Documents)) / float64(v)))
	}
}

// SplitText - split text
func SplitText(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return (r == '。' || r == '．' || r == '\n'/* || r == '！' || r == '？'*/)
	})
}

// ParseDocument split document to sentence and store 
func (d *Document) ParseDocument(text string) int {
	d.MorphemeCount = make(map[string]int)
	d.Plain = text
	for _, sentence := range SplitText(text) {
		s := Sentence{}
		s.ParseMeCab(sentence)
		for _, morp := range s.Morphemes {
			d.MorphemeCount[morp.Surface]++
		}
		d.Sentences = append(d.Sentences, s)
	}
	return 0
}

// AddDocument - Add document to Documents struct's slice
func (ds *Documents) AddDocument(text string) {
	d := Document{}
	d.ParseDocument(text)
	ds.Documents = append(ds.Documents, d)
}


// Vector calc functions
func InnerProd(sen1 Sentence, sen2 Sentence) float32 {
	sum := float32(0)
	for _, morp := range sen1.Morphemes {
		if (sen2.MorphemeMap[morp.Surface] != nil) {
			sum += morp.Value * sen2.MorphemeMap[morp.Surface].Value
		}
	}
	return sum
}
func VecAbs(sen Sentence) float32 {
	sum := float32(0)
	for _, morp := range sen.Morphemes {
		sum += morp.Value * morp.Value
	}
	return float32(math.Sqrt(float64(sum)))
}
