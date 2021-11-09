package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	P2CTable       map[string]([]string)
	CharProb       map[string]float64
	C2PTable       map[string](map[string]float64)
	WordsFreq      map[string]float64
	CharFreq       map[string]float64
	PinyinFreq     map[string]float64
	WordsFreq_Dim3 map[string]float64
)

func load_p2ctable() error {
	P2CTable = make(map[string][]string)
	data, _ := ioutil.ReadFile("./saved_model/P2CTable.json")
	json.Unmarshal(data, &P2CTable)
	return nil
}

func load_wordsfreq() {
	WordsFreq = make(map[string]float64)
	data, _ := ioutil.ReadFile("./saved_model/WordsFreq.json")
	json.Unmarshal(data, &WordsFreq)
}

func load_wordsfreq_dim3() {
	WordsFreq_Dim3 = make(map[string]float64)
	data, _ := ioutil.ReadFile("./saved_model/WordsFreq_dim3.json")
	json.Unmarshal(data, &WordsFreq_Dim3)
}

func load_pinyinfreq() {
	PinyinFreq = make(map[string]float64)
	data, _ := ioutil.ReadFile("./saved_model/PinyinFreq.json")
	json.Unmarshal(data, &PinyinFreq)
}

func load_charfreq() {
	CharFreq = make(map[string]float64)
	data, _ := ioutil.ReadFile("./saved_model/CharFreq.json")
	json.Unmarshal(data, &CharFreq)
}

func load_c2ptable() {
	C2PTable = make(map[string]map[string]float64)
	data, _ := ioutil.ReadFile("./saved_model/C2PTable.json")
	json.Unmarshal(data, &C2PTable)
}

func load_charprob() {
	CharProb = make(map[string]float64)
	data, _ := ioutil.ReadFile("./saved_model/CharProb.json")
	json.Unmarshal(data, &CharProb)
}

func LoadModel(dim int) {
	fmt.Println("Start Loading Model   Dim:", dim)
	load_p2ctable()
	load_c2ptable()
	load_charfreq()
	load_pinyinfreq()
	if dim == 3 {
		load_wordsfreq_dim3()
	}
	load_wordsfreq()
	load_charprob()
}
