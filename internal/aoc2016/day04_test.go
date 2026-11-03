package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestRoomDecryption(t *testing.T) {
	room := encryptedRoom{name: "qzmt-zixmtkozy-ivhz", sector: 343}
	if got := room.decryptedName(); got != "very encrypted name" {
		t.Fatalf("got %q, want very encrypted name", got)
	}
}

func Test_day04p01(t *testing.T) {
	aoc.AOCTest(t, day04p01, []aoc.TestInput{
		{Input: strings.NewReader("aaaaa-bbb-z-y-x-123[abxyz]\na-b-c-d-e-f-g-h-987[abcde]\nnot-a-real-room-404[oarel]\ntotally-real-room-200[decoy]"), Want: "1514"},
		{Input: aoc.FileInput(t, 2016, 4), Want: "137896"},
	})
}

func Test_day04p02(t *testing.T) {
	aoc.AOCTest(t, day04p02, []aoc.TestInput{
		{Input: strings.NewReader("northpole-object-storage-26[wrong]\nnorthpole-object-storage-0[oetra]"), Want: "0"},
		{Input: aoc.FileInput(t, 2016, 4), Want: "501"},
	})
}

func TestDay04ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day04p01", "northpole-object-storage-0[oetra]\n", day04p01},
		{"day04p02", "northpole-object-storage-0[oetra]\n", day04p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				reader := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(reader); !errors.Is(err, failure) {
					t.Fatalf("prefix %q: got %v, want input failure", prefix, err)
				}
			}
		})
	}
}

func TestDay04ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day04p01", day04p01},
		{"day04p02", day04p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay04MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"room", "room-xyz[abcde]", day04p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
