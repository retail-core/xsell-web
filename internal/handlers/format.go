package handlers

import (
	"fmt"
	"strconv"
	"strings"
)

func humanizeNumber(n float64) string {
	str := strconv.FormatFloat(n, 'f', -1, 64)
	parts := strings.Split(str, ".")
	intPart := parts[0]

	var result []string
	for i, digit := range reverse(intPart) {
		if i > 0 && i%3 == 0 {
			result = append(result, ",")
		}
		result = append(result, string(digit))
	}
	reversed := reverse(strings.Join(result, ""))

	if len(parts) > 1 {
		return fmt.Sprintf("%s.%s", reversed, parts[1])
	}
	return reversed
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}