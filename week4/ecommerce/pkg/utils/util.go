package utils

import "strconv"

func StrToInt(str string) int {
	if str == "" {
		return 0
	}

	result, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return result
}