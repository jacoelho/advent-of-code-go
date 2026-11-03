package aoc2021

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("D2FE28"), Want: "6"},
		{Input: strings.NewReader("8A004A801A8002F478"), Want: "16"},
		{Input: strings.NewReader("620080001611562C8802118E34"), Want: "12"},
		{Input: strings.NewReader("C0015000016115A2E0802F182340"), Want: "23"},
		{Input: strings.NewReader("A0016C880162017C3686B18A3D4780"), Want: "31"},
		{Input: aoc.FileInput(t, 2021, 16), Want: "949"},
	}

	aoc.AOCTest(t, day16p01, tests)
}

func Test_day16p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("D2FE28"), Want: "2021"},
		{Input: strings.NewReader("C200B40A82"), Want: "3"},
		{Input: strings.NewReader("04005AC33890"), Want: "54"},
		{Input: strings.NewReader("880086C3E88112"), Want: "7"},
		{Input: strings.NewReader("CE00C43D881120"), Want: "9"},
		{Input: strings.NewReader("D8005AC2A8F0"), Want: "1"},
		{Input: strings.NewReader("F600BC2D8F"), Want: "0"},
		{Input: strings.NewReader("9C005AC2F8F0"), Want: "0"},
		{Input: strings.NewReader("9C0141080250320F1802104A08"), Want: "1"},
		{Input: aoc.FileInput(t, 2021, 16), Want: "1114600142730"},
	}

	aoc.AOCTest(t, day16p02, tests)
}

func TestTransmissionErrors(t *testing.T) {
	for _, input := range []string{
		"", "D2FE2G", "D2FE2", "D2FE2F", "000000",
		"16004400",  // Type 5 comparison with one literal operand.
		"000028400", // Ten-bit payload containing an eleven-bit literal.
	} {
		t.Run(input, func(t *testing.T) {
			if _, err := day16p02(strings.NewReader(input)); err == nil {
				t.Fatal("expected invalid transmission error")
			}
		})
	}
}
