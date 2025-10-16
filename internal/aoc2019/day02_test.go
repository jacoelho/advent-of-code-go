package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/assert"
)

func Test_day02p01(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		program, err := ParseIntcodeInput(strings.NewReader(`1,9,10,3,2,3,11,0,99,30,40,50`))
		assert.NoError(t, err)

		result, err := runIntcodeWithInputs(program, 9, 10)
		assert.NoError(t, err)

		want := 3500
		assert.Equal(t, result, want)
	})

	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 2),
			Want:  "3706713",
		},
	}
	aoc.AOCTest(t, day2p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 2),
			Want:  "8609",
		},
	}
	aoc.AOCTest(t, day2p02, tests)
}
