package pinyin

import (
	"go-pinyin/utils"
	"log"
)

type HMM struct {
	Observation       []string
	Hidden            [][]string
	ProcessHiddenFunc func([]string) [][]string
	Pi                func(x, pinyin string) float64
	A                 func(x, y, pinyin string) float64
	B                 func(pinyin, x string) float64
}

func (h *HMM) Viterbi_dim2() []string {
	V := make([][]float64, 0)
	path := make(map[string](map[int][]string))

	Y := h.Observation
	X := h.Hidden

	col := make([]float64, 0)

	if len(X) <= 0 {
		return []string{}
	}

	for i := range X[0] {
		v := h.Pi(X[0][i], Y[0]) * h.B(Y[0], X[0][i])
		col = append(col, v)
		path[X[0][i]] = make(map[int][]string)
		path[X[0][i]][0] = []string{X[0][i]}
	}
	V = append(V, col)

	for t := 1; t < len(X); t++ {
		col = make([]float64, len(X[t]))
		for i := range col {
			col[i] = -1
		}

		V = append(V, col)

		for i, s := range X[t] {
			var prob float64 = -1
			for j, _s := range X[t-1] {
				nprob := V[t-1][j] * h.A(_s, s, Y[t]) * h.B(Y[t], s)
				if nprob > prob {
					prob = nprob
					V[t][i] = prob
					path[s] = make(map[int][]string)
					path[s][t] = make([]string, len(path[_s][t-1]))

					copy(path[s][t], path[_s][t-1])
					path[s][t] = append(path[s][t], s)
				}
			}
		}
	}
	var prob float64 = -1
	state := ""
	for i := range X[len(X)-1] {
		if V[len(X)-1][i] > prob {
			prob = V[len(X)-1][i]

			state = X[len(X)-1][i]
		}
	}
	return path[state][len(X)-1]
}

func (h *HMM) Viterbi_dim3() []string {
	V := make([][]float64, 0)
	path := make(map[string](map[int][]string))

	Y := h.Observation
	X := h.Hidden

	col := make([]float64, 0)

	if len(X) <= 0 {
		return []string{}
	}

	for i := range X[0] {
		var v float64
		if len(Y) == 1 {
			v = h.Pi(X[0][i], Y[0]) * h.B(Y[0], X[0][i])
		} else {
			x0 := X[0][i][0:3]
			x1 := X[0][i][3:]
			v = h.Pi(X[0][i], Y[0]) * h.B(Y[0], x0) * h.B(Y[1], x1)
		}
		col = append(col, v)
		path[X[0][i]] = make(map[int][]string)
		path[X[0][i]][0] = []string{X[0][i]}
	}
	V = append(V, col)

	for t := 1; t < len(X); t++ {
		col = make([]float64, len(X[t]))
		for i := range col {
			col[i] = -1
		}

		V = append(V, col)

		for i, s := range X[t] {
			var prob float64 = -1
			for j, _s := range X[t-1] {
				if _s[3:] != s[0:3] {
					continue
				}
				nprob := V[t-1][j] * h.A(_s, s, Y[t+1]) * h.B(Y[t+1], s[3:])
				if nprob > prob {
					prob = nprob

					V[t][i] = prob
					path[s] = make(map[int][]string)
					path[s][t] = make([]string, len(path[_s][t-1]))

					copy(path[s][t], path[_s][t-1])
					path[s][t] = append(path[s][t], s)
				}
			}
		}
	}
	var prob float64 = -1
	state := ""
	for i := range X[len(X)-1] {
		if V[len(X)-1][i] > prob {
			prob = V[len(X)-1][i]

			state = X[len(X)-1][i]
		}
	}
	ans := path[state][len(X)-1]
	if len(ans) == 0 {
		return ans
	}
	output := make([]string, 0)
	output = append(output, ans[0])
	for i := 1; i < len(ans); i++ {
		output = append(output, ans[i][3:])
	}
	return output

}

func (h *HMM) Run(obs []string, dim int) string {
	h.Observation = obs
	h.Hidden = h.ProcessHiddenFunc(obs)
	var ans []string
	if dim == 2 {
		ans = h.Viterbi_dim2()
	} else if dim == 3 {
		ans = h.Viterbi_dim3()
	} else {
		log.Fatalf("Error: Unimplemented Dim: %d\n", dim)
		return ""
	}
	ret := ""
	for _, s := range ans {
		ret = ret + s
	}
	return ret
}

func A(x, y, pinyin string) float64 {
	lambda := 0.0001
	if v, ok := utils.WordsFreq[x+y]; ok {
		return v / utils.CharFreq[x]
	}
	if v, ok := utils.CharProb[y+"-"+pinyin]; ok {
		return float64(lambda) * v
	}
	return 1e-10
}

func A_3(x, y, pinyin string) float64 {
	gamma := 0.01
	if x[3:] != y[0:3] {
		return 0
	}
	if v, ok := utils.WordsFreq_Dim3[x+y[3:]]; ok {
		return v / utils.WordsFreq[x]
	} else {
		return gamma * A(y[0:3], y[3:], pinyin)
	}
}

func B(pinyin, x string) float64 {
	if _, ok := utils.C2PTable[x]; !ok {
		return 1
	}
	if v, ok := (utils.C2PTable[x])[pinyin]; ok {
		return v / utils.C2PTable[x]["sum"]
	}
	return 1 / utils.C2PTable[x]["sum"]
}

func Pi(x, pinyin string) float64 {
	if v, ok := utils.CharProb[x+"-"+pinyin]; ok {
		return v
	}
	return 1e-10
}

func Pi_3(x, pinyin string) float64 {
	if len(x) <= 3 {
		return Pi(x, pinyin)
	}
	if v, ok := utils.WordsFreq[x]; ok {
		return v
	}
	return 1
}
