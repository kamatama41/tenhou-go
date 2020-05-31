package tenhou

import (
	"strconv"
	"strings"
)

func splitByCommaAsInt(str string) ([]int, error) {
	var res []int
	for _, s := range splitByComma(str) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func splitByCommaAsHaiList(str string) (HaiList, error) {
	var res HaiList
	for _, s := range splitByComma(str) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		hai := Hai{}
		if err := hai.Unmarshal(i); err != nil {
			return nil, err
		}
		res = append(res, hai)
	}
	return res, nil
}

func splitByComma(str string) []string {
	return strings.Split(str, ",")
}

func joinByCommaFromInts(ints []int) string {
	var res []string
	for _, i := range ints {
		res = append(res, strconv.Itoa(i))
	}
	return joinByComma(res)
}

func joinByComma(strs []string) string {
	return strings.Join(strs, ",")
}
