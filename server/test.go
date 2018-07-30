package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"bufio"

	"./ceft"
)

func printResult(sen ceft.Sentence) {
	fmt.Println(sen.Plain)
	fmt.Println("	解析結果: ")
	for _, ele := range sen.Morphemes {
		fmt.Printf("		%s:	TF: %f	自立語->%v\n	%v\n", ele.Surface, sen.CalcTF(ele), ele.IsIndep(), ele.Feature)
	}
}

func main() {
	sen := ceft.Sentence{}
	Docs := ceft.Documents{}
	str := "今日も部室のがっこうぐらし！を何冊も読んだ．"
	for index, arg := range os.Args {
		if (index == 1) {
			str = arg
		} else if (index > 1) {
			if (arg[len(arg) - 5:] == ".json") {
				fp, err := os.Open(arg)
				if (err == nil) {
					scanner := bufio.NewScanner(fp)
					for scanner.Scan() {
						docj := ceft.DecodeDocJSON(scanner.Text())
						Docs.AddDocument(docj.Text)
						fmt.Printf("\033[0G\033[K(%d) Loaded: %s", len(Docs.Documents), docj.Title)
						if err := scanner.Err(); err != nil {
							fmt.Println()
							fmt.Println(err)
						}
					}
					fp.Close()
					fmt.Println("")
				}
			} else {
				bts, err := ioutil.ReadFile(arg)
				if (err == nil) {
					Docs.AddDocument(string(bts))
				}
			}
		}
	}
	
	if (sen.ParseMeCab(str) == 0) {
		printResult(sen)
		
		if (len(Docs.Documents) > 1) {
			Docs.CalcIDF()

			max := float32(0)
			ressen := ceft.Sentence{}
			for didx, _ := range Docs.Documents {
				for sidx, _ := range Docs.Documents[didx].Sentences {
					for midx, m := range Docs.Documents[didx].Sentences[sidx].Morphemes {
						Docs.Documents[didx].Sentences[sidx].Morphemes[midx].Value = Docs.Documents[didx].Sentences[sidx].CalcTF(m) * Docs.IDF[m.Surface]
					}
					tmp := ceft.InnerProd(sen, Docs.Documents[didx].Sentences[sidx]) / (ceft.VecAbs(sen) * ceft.VecAbs(Docs.Documents[didx].Sentences[sidx]))
					if (tmp > max) {
						max = tmp
						ressen = Docs.Documents[didx].Sentences[sidx]
					}
				}
			}
			fmt.Println("cos: ", max)
			printResult(ressen)
		}
	}
}
