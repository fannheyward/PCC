package main

import "strconv"

func errInfo(info string) map[string]interface{} {
	return map[string]interface{}{
		"error_message": info,
	}
}

func strToInt(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}
