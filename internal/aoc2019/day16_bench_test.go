package aoc2019

import (
	"bytes"
	"io"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func BenchmarkDay16Part1(b *testing.B) {
	input, err := io.ReadAll(aoc.FileInput(b, 2019, 16))
	if err != nil {
		b.Fatal(err)
	}
	for _, fixture := range []struct {
		name  string
		input []byte
		want  string
	}{
		{"sample", []byte("80871224585914546619083218645595"), "24176176"},
		{"input", input, "40921727"},
	} {
		b.Run(fixture.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				got, err := day16p01(bytes.NewReader(fixture.input))
				if err != nil || got != fixture.want {
					b.Fatalf("got %q, %v; want %q", got, err, fixture.want)
				}
			}
		})
	}
}
