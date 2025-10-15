package aoc2024

import (
	"io"
	"iter"
	"maps"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/conc"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseSecretNumbers(r io.Reader) ([]int, error) {
	s := scanner.NewScanner[int](r, convert.ScanNumber)
	return slices.Collect(s.Values()), s.Err()
}

func secretNumberIter(n int) iter.Seq[int] {
	return xiter.Apply(n, func(i int) int {
		n = (n ^ (n << 6)) & 16777215
		n = (n ^ (n >> 5)) & 16777215
		n = (n ^ (n << 11)) & 16777215
		return n
	})
}

func secretNumberIterationN(n int) int {
	v, ok := xiter.Next(xiter.Skip(secretNumberIter(n), 2000))
	if !ok {
		panic("next not found")
	}
	return v
}

func day22p01(r io.Reader) (string, error) {
	secretNumbers := aoc.Must(parseSecretNumbers(r))

	var total int
	for _, secretNumber := range secretNumbers {
		total += secretNumberIterationN(secretNumber)
	}
	return strconv.Itoa(total), nil
}

func sellSequences(n int) map[[4]int]int {
	digitSecretNumberIter := xiter.Map(func(v int) int { return v % 10 }, xiter.Take(2000, secretNumberIter(n)))

	seen := collections.NewSet[[4]int]()
	total := make(map[[4]int]int)
	for v := range xiter.Window(5, digitSecretNumberIter) {
		a, b, c, d, e := v[0], v[1], v[2], v[3], v[4]
		seq := [4]int{b - a, c - b, d - c, e - d}
		if seen.Contains(seq) {
			continue
		}
		seen.Add(seq)
		total[seq] = total[seq] + e
	}
	return total
}

func day22p02(r io.Reader) (string, error) {
	secretNumbers := aoc.Must(parseSecretNumbers(r))

	total := make(map[[4]int]int)
	for _, m := range conc.Map(sellSequences, secretNumbers) {
		for k, v := range m {
			total[k] += v
		}
	}
	maxBananas := xiter.Max(maps.Values(total))
	return strconv.Itoa(maxBananas), nil
}
