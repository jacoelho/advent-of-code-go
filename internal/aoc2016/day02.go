package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func bathroomCode(r io.Reader, keypad []string, position grid.Position2D[int]) (string, error) {
	var code strings.Builder
	s := bufio.NewScanner(r)
	for s.Scan() {
		for _, direction := range s.Bytes() {
			next := position
			switch direction {
			case 'U':
				next.Y--
			case 'D':
				next.Y++
			case 'L':
				next.X--
			case 'R':
				next.X++
			default:
				return "", fmt.Errorf("invalid keypad direction %q", direction)
			}
			if next.Y >= 0 && next.Y < len(keypad) && next.X >= 0 && next.X < len(keypad[next.Y]) && keypad[next.Y][next.X] != ' ' {
				position = next
			}
		}
		code.WriteByte(keypad[position.Y][position.X])
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	return code.String(), nil
}

func day02p01(r io.Reader) (string, error) {
	return bathroomCode(r, []string{"123", "456", "789"}, grid.Position2D[int]{X: 1, Y: 1})
}

func day02p02(r io.Reader) (string, error) {
	return bathroomCode(r, []string{"  1  ", " 234 ", "56789", " ABC ", "  D  "}, grid.Position2D[int]{Y: 2})
}
