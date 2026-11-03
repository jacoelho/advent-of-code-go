package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day07Example = `$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k`

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day07Example), Want: "95437"},
		{Input: aoc.FileInput(t, 2022, 7), Want: "1427048"},
	}
	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day07Example), Want: "24933642"},
		{Input: aoc.FileInput(t, 2022, 7), Want: "2940614"},
	}
	aoc.AOCTest(t, day07p02, tests)
}

func TestDirectoryIdentity(t *testing.T) {
	const input = "$ cd /\n$ ls\ndir a\ndir empty\n10 root\n$ cd a\n$ ls\n5 child\n$ cd /\n$ ls\n10 root"
	aoc.AOCTest(t, day07p01, []aoc.TestInput{{Input: strings.NewReader(input), Want: "20"}})
	aoc.AOCTest(t, day07p02, []aoc.TestInput{{Input: strings.NewReader(input), Want: "0"}})
	if _, err := day07p01(strings.NewReader("$ ls\n1 a\n2 a")); err == nil {
		t.Fatal("expected conflicting listing error")
	}
}

func TestDay07ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day07p01},
		{"part2", day07p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
