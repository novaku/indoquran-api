package helpers

import (
	"strconv"
	"strings"
)

// StringInSlice : array check string
func StringInSlice(search string, list []string) bool {
	for _, str := range list {
		if str == search {
			return true
		}
	}
	return false
}

// ReverseString : reverse string
func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// EngNumberToArNumber : english number to arab number (rtl)
func EngNumberToArNumber(no int) string {
	str := ReverseString(strconv.Itoa(no))

	replace := map[string]string{
		"0": "٠",
		"1": "١",
		"2": "٢",
		"3": "٣",
		"4": "٤",
		"5": "٥",
		"6": "٦",
		"7": "٧",
		"8": "٨",
		"9": "٩",
	}

	for s, r := range replace {
		str = strings.Replace(str, s, r, -1)
	}

	return str
}
