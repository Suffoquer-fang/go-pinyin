package main

import (
	"bufio"
	"fmt"
	"go-pinyin/pinyin"
	"go-pinyin/utils"
	"io"
	"os"
	"strings"
)

func HandleFiles(dim int, model *pinyin.HMM, inpath, outpath string) {
	f, _ := os.Open(inpath)
	buf := bufio.NewReader(f)
	out_file, _ := os.Create(outpath)
	var i int = 1
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			return
		}
		line = strings.TrimSpace(line)
		obs := utils.ProcessObs(line)

		ans := model.Run(obs, dim)
		io.WriteString(out_file, ans+"\n")
		fmt.Println("Done with line", i)
		i = i + 1
	}
}

func CalcAccuracy(outpath, std_path string) {
	out_file, _ := os.Open(outpath)
	std_file, _ := os.Open(std_path)
	out_data, _ := io.ReadAll(out_file)
	std_data, _ := io.ReadAll(std_file)

	out_data_lines := strings.Split(string(out_data), "\n")
	std_data_lines := strings.Split(string(std_data), "\r\n")

	tot_word := 0
	correct_word := 0

	tot_sen := 0
	correct_sen := 0
	for i := range out_data_lines {
		if len(out_data_lines[i]) == 0 {
			continue
		}
		tot_sen += 1
		if out_data_lines[i] == std_data_lines[i] {
			correct_sen += 1
		}

		for j := 0; j < len(out_data_lines[i]); j += 3 {
			if out_data_lines[i][j:j+3] == std_data_lines[i][j:j+3] {
				correct_word += 1
			}
			tot_word += 1
		}
	}
	fmt.Println("tot_word =", tot_word)
	fmt.Println("correct_word =", correct_word)
	fmt.Println("tot_sen =", tot_sen)
	fmt.Println("correct_sen =", correct_sen)
	fmt.Printf("WordAccuracy: %.2f\n", float64(100.0*correct_word/tot_word))
	fmt.Printf("SentenceAccuracy: %.2f\n", float64(100*correct_sen/tot_sen))
}

func main() {
	dim := 2
	utils.LoadModel(dim)
	var h *pinyin.HMM
	if dim == 2 {
		h = &pinyin.HMM{
			ProcessHiddenFunc: utils.ProcessHidden,
			Pi:                pinyin.Pi,
			A:                 pinyin.A,
			B:                 pinyin.B,
		}
	} else if dim == 3 {
		h = &pinyin.HMM{
			ProcessHiddenFunc: utils.ProcessHidden_3,
			Pi:                pinyin.Pi_3,
			A:                 pinyin.A_3,
			B:                 pinyin.B,
		}
	} else {
		h = nil
	}

	fmt.Println("Model Loaded")
	// HandleFiles(dim, h, "input.txt", "output.txt")
	// CalcAccuracy("output.txt", "std_output.txt")

	input_pinyin := ""
	reader := bufio.NewReader(os.Stdin)
	for {
		input_pinyin, _ = reader.ReadString('\n')
		obs := utils.ProcessObs(input_pinyin)
		ans := h.Run(obs, dim)
		fmt.Println(ans)
	}
}
