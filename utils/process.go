package utils

import "strings"

func ProcessHidden(obs []string) []([]string) {
	ret := make([]([]string), 0)
	for _, p := range obs {
		if v, ok := P2CTable[string(p)]; ok {
			ret = append(ret, v)
		}
	}
	return ret
}

func ProcessHidden_3(obs []string) [][]string {
	ret := make([][]string, 0)
	if len(obs) == 1 {
		return [][]string{P2CTable[obs[0]]}
	}
	for i := 0; i < len(obs)-1; i++ {
		ls := make([]string, 0)
		if v, ok := P2CTable[obs[i]]; ok {
			if _v, _ok := P2CTable[obs[i+1]]; _ok {
				for _, w := range v {
					for _, _w := range _v {
						ls = append(ls, string(w+_w))
					}
				}
			}
		}
		if len(ls) > 0 {
			ret = append(ret, ls)
		}
	}
	return ret
}

func ProcessObs(raw_obs string) []string {
	raw_obs = strings.Trim(raw_obs, " \n")
	obs := strings.Split(raw_obs, " ")
	ret_obs := make([]string, 0)
	for _, p := range obs {
		if _, ok := P2CTable[p]; ok {
			ret_obs = append(ret_obs, p)
		}
	}
	return ret_obs
}
