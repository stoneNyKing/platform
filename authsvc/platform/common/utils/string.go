package utils

import "strings"

/*
	去掉数字字符串中的逗号，比如
	 1,711,624,950 变成 1711624950
*/

func ClearComma(s string) string {
	if s == "" {
		return ""
	}
	sret := ""
	ss := strings.Split(s,",")

	for _,sv := range ss {
		sret += sv
	}

	return sret
}

/*
	去掉数字字符串的%,比如：
	3.65% 变为 3.65
*/

func ClearPercent(s string) string {
	if s == "" {
		return s
	}

	ss := strings.Split(s,"%")

	return ss[0]
}