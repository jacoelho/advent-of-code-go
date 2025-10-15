package aoc2024

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseStoneArrangement(r io.Reader) []int64 {
	var result []int64
	b := bufio.NewScanner(r)
	for b.Scan() {
		for v := range strings.FieldsSeq(b.Text()) {
			result = append(result, aoc.Must(strconv.ParseInt(v, 10, 64)))
		}
	}
	return result
}

func memoizeCount() func(int64, int64) int64 {
	cache := make(map[[2]int64]int64)

	var countRec func(stone int64, steps int64) int64
	countRec = func(stone int64, steps int64) int64 {
		if steps == 0 {
			return 1
		}

		k := [2]int64{stone, steps}
		if v, ok := cache[k]; ok {
			return v
		}

		if stone == 0 {
			cache[k] = countRec(1, steps-1)
			return cache[k]
		}
		digits := convert.ToDigits(stone)
		if len(digits)%2 != 0 {
			cache[k] = countRec(stone*2024, steps-1)
			return cache[k]
		}

		mid := len(digits) / 2
		n1 := convert.FromDigits(digits[:mid])
		n2 := convert.FromDigits(digits[mid:])
		cache[k] = countRec(n1, steps-1) + countRec(n2, steps-1)
		return cache[k]
	}

	return countRec
}

func countBlinks(stones []int64, steps int64) int64 {
	countFunc := memoizeCount()

	return xslices.Reduce(func(sum int64, stone int64) int64 {
		return sum + countFunc(stone, steps)
	}, 0, stones)
}

func day11p01(r io.Reader) (string, error) {
	stones := parseStoneArrangement(r)

	total := countBlinks(stones, 25)

	return strconv.FormatInt(total, 10), nil
}

func day11p02(r io.Reader) (string, error) {
	stones := parseStoneArrangement(r)

	total := countBlinks(stones, 75)

	return strconv.FormatInt(total, 10), nil
}
