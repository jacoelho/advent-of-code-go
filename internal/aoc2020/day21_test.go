package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day21p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`mxmxvkd kfcds sqjhc nhms (contains dairy, fish)
trh fvjkl sbzzf mxmxvkd (contains dairy)
sqjhc fvjkl (contains soy)
sqjhc mxmxvkd sbzzf (contains fish)`),
			Want: "5",
		},
		{
			Input: aoc.FileInput(t, 2020, 21),
			Want:  "2389",
		},
	}

	aoc.AOCTest(t, day21p01, tests)
}

func Test_day21p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`mxmxvkd kfcds sqjhc nhms (contains dairy, fish)
trh fvjkl sbzzf mxmxvkd (contains dairy)
sqjhc fvjkl (contains soy)
sqjhc mxmxvkd sbzzf (contains fish)`),
			Want: "mxmxvkd,sqjhc,fvjkl",
		},
		{
			Input: aoc.FileInput(t, 2020, 21),
			Want:  "fsr,skrxt,lqbcg,mgbv,dvjrrkv,ndnlm,xcljh,zbhp",
		},
	}

	aoc.AOCTest(t, day21p02, tests)
}
