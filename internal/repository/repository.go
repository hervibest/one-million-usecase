package repository

import "fmt"

func generateQuestionsMark(n int) []string {
	s := make([]string, n)
	for i := range s {
		s[i] = fmt.Sprintf("$%d", i+1) // Menggunakan placeholder PostgreSQL ($1, $2, dst.)
	}
	return s
}
