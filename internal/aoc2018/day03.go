package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type fabricClaim struct {
	id, x, y, width, height int
}

func claimedFabric(r io.Reader) ([]fabricClaim, map[grid.Position2D[int]]int, error) {
	s := bufio.NewScanner(r)
	var claims []fabricClaim
	fabric := make(map[grid.Position2D[int]]int)
	ids := make(map[int]bool)
	for s.Scan() {
		var claim fabricClaim
		if _, err := fmt.Sscanf(s.Text(), "#%d @ %d,%d: %dx%d", &claim.id, &claim.x, &claim.y, &claim.width, &claim.height); err != nil {
			return nil, nil, fmt.Errorf("invalid claim %q: %w", s.Text(), err)
		}
		if claim.id <= 0 || claim.x < 0 || claim.y < 0 || claim.width <= 0 || claim.height <= 0 || ids[claim.id] {
			return nil, nil, fmt.Errorf("invalid or duplicate claim %q", s.Text())
		}
		ids[claim.id] = true
		claims = append(claims, claim)
		for y := range claim.height {
			for x := range claim.width {
				fabric[grid.Position2D[int]{X: claim.x + x, Y: claim.y + y}]++
			}
		}
	}
	return claims, fabric, s.Err()
}

func day03p01(r io.Reader) (string, error) {
	_, fabric, err := claimedFabric(r)
	if err != nil {
		return "", err
	}

	overlap := 0
	for _, claims := range fabric {
		if claims > 1 {
			overlap++
		}
	}
	return strconv.Itoa(overlap), nil
}

func day03p02(r io.Reader) (string, error) {
	claims, fabric, err := claimedFabric(r)
	if err != nil {
		return "", err
	}

	for _, claim := range claims {
		intact := true
		for y := range claim.height {
			for x := range claim.width {
				if fabric[grid.Position2D[int]{X: claim.x + x, Y: claim.y + y}] != 1 {
					intact = false
				}
			}
		}
		if intact {
			return strconv.Itoa(claim.id), nil
		}
	}
	return "", fmt.Errorf("no intact fabric claim")
}
