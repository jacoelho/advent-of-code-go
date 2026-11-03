package aoc2016

import (
	"strings"
	"testing"

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
