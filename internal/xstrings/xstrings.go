package xstrings

func SubSlices(s string, n int) [][2]string {
	if n > len(s) {
		n = len(s)
	}
	result := make([][2]string, n)
	for i := 1; i <= n; i++ {
		result[i-1] = [2]string{s[:i], s[i:]}
	}
	return result
}
