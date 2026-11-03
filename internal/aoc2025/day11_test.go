package aoc2025

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out`),
			Want: "5",
		},
		{
			Input: aoc.FileInput(t, 2025, 11),
			Want:  "696",
		},
	}
	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out`),
			Want: "2",
		},
		{
			Input: aoc.FileInput(t, 2025, 11),
			Want:  "473741288064360",
		},
	}
	aoc.AOCTest(t, day11p02, tests)
}

func TestDevicePathFailures(t *testing.T) {
	var overflowing strings.Builder
	overflowing.WriteString("you: n0\nsvr: dac\ndac: fft\nfft: n0\n")
	for i := range strconv.IntSize - 1 {
		next := fmt.Sprintf("n%d", i+1)
		if i == strconv.IntSize-2 {
			next = "out"
		}
		fmt.Fprintf(&overflowing, "n%d: %s %s\n", i, next, next)
	}
	for _, test := range []struct{ name, input string }{
		{"empty list", ""},
		{"empty record", "you: out\n \n"},
		{"reachable cycle", "you: a\nsvr: dac\ndac: fft\nfft: a\na: b\nb: a out"},
		{"path count overflow", overflowing.String()},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, solve := range []func(io.Reader) (string, error){day11p01, day11p02} {
				if got, err := solve(strings.NewReader(test.input)); err == nil || got != "" {
					t.Fatalf("got %q, %v; want an error without a result", got, err)
				}
			}
		})
	}
}

func TestDeviceTerminalPaths(t *testing.T) {
	input := "you: out out\nsvr: dac\ndac: fft\nfft: out out\nout: you"
	for _, solve := range []func(io.Reader) (string, error){day11p01, day11p02} {
		got, err := solve(strings.NewReader(input))
		if err != nil || got != "2" {
			t.Fatalf("got %q, %v; want two parallel paths ending at out", got, err)
		}
	}
}

func TestDeviceLabels(t *testing.T) {
	for _, test := range []struct {
		name, input string
		want        [2]string
	}{
		{
			name:  "forward references and required visits",
			input: "you: merge\nsvr: dac fft both out\nboth: dac\ndac: fft merge\nfft: merge\nmerge: out sink",
			want:  [2]string{"1", "2"},
		},
		{
			name:  "zero contribution target is terminal",
			input: "you: out\nsvr: out\nout: svr",
			want:  [2]string{"1", "0"},
		},
		{
			name:  "last definition replaces edges",
			input: "you: out\nyou: out out\nsvr: dac\ndac: fft\nfft: out",
			want:  [2]string{"2", "1"},
		},
		{
			name:  "missing starting devices",
			input: "other: out",
			want:  [2]string{"0", "0"},
		},
		{
			name:  "mixed separators and colonless records",
			input: ": :you :: out : out:\nsvr dac\n::dac: fft:\nfft out",
			want:  [2]string{"2", "1"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for part, solve := range []func(io.Reader) (string, error){day11p01, day11p02} {
				got, err := solve(strings.NewReader(test.input))
				if err != nil || got != test.want[part] {
					t.Fatalf("part %d: got %q, %v; want %q", part+1, got, err, test.want[part])
				}
			}
		})
	}
}

func TestDeviceRequiredTarget(t *testing.T) {
	network, err := parseDevicesList(strings.NewReader("source: out out"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := countPaths(network, []string{"out"}, "out", "source")
	if err != nil || got != 2 {
		t.Fatalf("got %d, %v; want two paths through the required target", got, err)
	}
}

func TestDay11ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, solve := range []func(io.Reader) (string, error){day11p01, day11p02} {
		for _, prefix := range []string{"", "you: out\n"} {
			input := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
			if _, err := solve(input); !errors.Is(err, failure) {
				t.Fatalf("prefix %q: got %v; want reader failure", prefix, err)
			}
		}
	}
}
