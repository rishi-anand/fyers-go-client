package utils

import "fmt"

func FormatStrArrWithQuotes(in []string) []string {
	for i := 0; i < len(in); i++ {
		in[i] = FormatStrWithQuotes(in[i])
	}
	return in
}
func FormatStrWithQuotes(in string) string {
	if in[0] == '"' && in[len(in)-1] == '"' {
		return in
	}
	return fmt.Sprintf("\"%s\"", in)
}
