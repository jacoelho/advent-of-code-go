package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type fabricClaim struct {
	id, x, y, width, height int
}

type fabric struct {
	width  int
	claims []int
}

func (f fabric) intact(claim fabricClaim) bool {
	for y := claim.y; y < claim.y+claim.height; y++ {
		for _, claims := range f.claims[y*f.width+claim.x : y*f.width+claim.x+claim.width] {
			if claims != 1 {
				return false
			}
		}
	}
	return true
}

func claimedFabric(r io.Reader) ([]fabricClaim, fabric, error) {
	s := bufio.NewScanner(r)
	var claims []fabricClaim
	ids := make(map[int]bool)
	var width, height int
	for s.Scan() {
		// int32 fields keep every coordinate sum and the fabric area within int.
		var id, x, y, w, h int32
		if _, err := fmt.Sscanf(s.Text(), "#%d @ %d,%d: %dx%d", &id, &x, &y, &w, &h); err != nil {
			return nil, fabric{}, fmt.Errorf("invalid claim %q: %w", s.Text(), err)
		}
		claim := fabricClaim{id: int(id), x: int(x), y: int(y), width: int(w), height: int(h)}
		if claim.id <= 0 || claim.x < 0 || claim.y < 0 || claim.width <= 0 || claim.height <= 0 || ids[claim.id] {
			return nil, fabric{}, fmt.Errorf("invalid or duplicate claim %q", s.Text())
		}
		ids[claim.id] = true
		claims = append(claims, claim)
		width = max(width, claim.x+claim.width)
		height = max(height, claim.y+claim.height)
	}
	if err := s.Err(); err != nil {
		return nil, fabric{}, err
	}
	if len(claims) == 0 {
		return nil, fabric{}, fmt.Errorf("no fabric claims")
	}

	f := fabric{width: width, claims: make([]int, width*height)}
	for _, claim := range claims {
		for y := claim.y; y < claim.y+claim.height; y++ {
			for x := y*width + claim.x; x < y*width+claim.x+claim.width; x++ {
				f.claims[x]++
			}
		}
	}
	return claims, f, nil
}

func day03p01(r io.Reader) (string, error) {
	_, f, err := claimedFabric(r)
	if err != nil {
		return "", err
	}

	overlap := 0
	for _, claims := range f.claims {
		if claims > 1 {
			overlap++
		}
	}
	return strconv.Itoa(overlap), nil
}

func day03p02(r io.Reader) (string, error) {
	claims, f, err := claimedFabric(r)
	if err != nil {
		return "", err
	}

	for _, claim := range claims {
		if f.intact(claim) {
			return strconv.Itoa(claim.id), nil
		}
	}
	return "", fmt.Errorf("no intact fabric claim")
}
