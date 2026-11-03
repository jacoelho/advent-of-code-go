package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	aoc.AOCTest(t, day20p01, []aoc.TestInput{
		{Input: strings.NewReader("p=<3,0,0>, v=<2,0,0>, a=<-1,0,0>\np=<4,0,0>, v=<0,0,0>, a=<-2,0,0>"), Want: "0"},
		{Input: strings.NewReader("p=<0,0,0>, v=<2,0,0>, a=<1,0,0>\np=<1,0,0>, v=<-2,0,0>, a=<1,0,0>"), Want: "1"},
		{Input: strings.NewReader("p=<0,0,0>, v=<0,0,0>, a=<-9223372036854775808,0,0>\np=<0,0,0>, v=<0,0,0>, a=<9223372036854775807,0,0>"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 20), Want: "243"},
	})
}

func Test_day20p02(t *testing.T) {
	aoc.AOCTest(t, day20p02, []aoc.TestInput{
		{Input: strings.NewReader("p=<-6,0,0>, v=<3,0,0>, a=<0,0,0>\np=<-4,0,0>, v=<2,0,0>, a=<0,0,0>\np=<-2,0,0>, v=<1,0,0>, a=<0,0,0>\np=<3,0,0>, v=<-1,0,0>, a=<0,0,0>"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 20), Want: "648"},
	})
}

func TestParticleCollisions(t *testing.T) {
	tests := []struct{ name, input, want string }{
		{
			name:  "after two thousand ticks",
			input: "p=<0,0,0>, v=<1,0,0>, a=<0,0,0>\np=<4002,0,0>, v=<-1,0,0>, a=<0,0,0>",
			want:  "0",
		},
		{
			name:  "fractional crossing",
			input: "p=<0,0,0>, v=<2,0,0>, a=<0,0,0>\np=<3,0,0>, v=<0,0,0>, a=<0,0,0>",
			want:  "2",
		},
		{
			name:  "initial collision",
			input: "p=<0,0,0>, v=<1,0,0>, a=<0,0,0>\np=<0,0,0>, v=<-1,0,0>, a=<0,0,0>",
			want:  "0",
		},
		{
			name:  "simultaneous removal before later intersection",
			input: "p=<-2,0,0>, v=<1,0,0>, a=<0,0,0>\np=<0,0,0>, v=<0,0,0>, a=<0,0,0>\np=<2,0,0>, v=<-1,0,0>, a=<0,0,0>\np=<3,0,0>, v=<-1,0,0>, a=<0,0,0>",
			want:  "1",
		},
		{
			name:  "identical trajectories",
			input: "p=<1,2,3>, v=<4,5,6>, a=<7,8,9>\np=<1,2,3>, v=<4,5,6>, a=<7,8,9>",
			want:  "0",
		},
		{
			name:  "quadratic collision",
			input: "p=<0,0,0>, v=<0,0,0>, a=<1,0,0>\np=<6,0,0>, v=<0,0,0>, a=<0,0,0>",
			want:  "0",
		},
		{
			name:  "collision time exceeds int64",
			input: "p=<-9223372036854775808,0,0>, v=<1,0,0>, a=<0,0,0>\np=<9223372036854775807,0,0>, v=<0,0,0>, a=<0,0,0>",
			want:  "0",
		},
		{
			name:  "different axis times",
			input: "p=<0,0,0>, v=<1,1,0>, a=<0,0,0>\np=<2,3,0>, v=<0,0,0>, a=<0,0,0>",
			want:  "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := day20p02(strings.NewReader(tt.input))
			if err != nil || got != tt.want {
				t.Fatalf("day20p02() = %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}
