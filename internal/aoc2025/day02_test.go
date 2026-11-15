package aoc2025

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`),
			Want:  "1227775554",
		},
		{
			Input: aoc.FileInput(t, 2025, 02),
			Want:  "24157613387",
		},
	}
	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`),
			Want:  "4174379265",
		},
		{
			Input: aoc.FileInput(t, 2025, 02),
			Want:  "33832678380",
		},
	}
	aoc.AOCTest(t, day02p02, tests)
}

func TestRepeatedBlockSumMatchesDigitComparison(t *testing.T) {
	isRepeated := func(id int, repeats func(int) bool) bool {
		digits := strconv.Itoa(id)
		for block := 1; block <= len(digits)/2; block++ {
			if len(digits)%block == 0 && repeats(len(digits)/block) &&
				strings.Repeat(digits[:block], len(digits)/block) == digits {
				return true
			}
		}
		return false
	}
	random := rand.New(rand.NewPCG(7, 8))
	for _, repeats := range []func(int) bool{
		func(times int) bool { return times == 2 },
		func(times int) bool { return times >= 2 },
	} {
		for range 300 {
			start := random.IntN(2_000_000)
			p := productRange{start: start, end: start + random.IntN(5_000)}
			want := 0
			for id := p.start; id <= p.end; id++ {
				if isRepeated(id, repeats) {
					want += id
				}
			}
			if got, err := repeatedBlockSum(p, repeats); err != nil || got != want {
				t.Fatalf("range %v: got %d, %v; want %d", p, got, err, want)
			}
		}
	}
}
