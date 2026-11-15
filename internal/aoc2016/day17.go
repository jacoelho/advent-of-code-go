package aoc2016

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type vaultPath struct {
	position grid.Position2D[int]
	path     string
}

func (p vaultPath) atVault() bool {
	return p.position == (grid.Position2D[int]{X: 3, Y: 3})
}

func vaultPasscode(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.ContainsAny(data, " \t\r\n") {
		return "", fmt.Errorf("invalid vault passcode %q", data)
	}
	return string(data), nil
}

func vaultNeighbours(passcode string) func(vaultPath) []vaultPath {
	return func(current vaultPath) []vaultPath {
		if current.atVault() {
			return nil
		}
		digest := md5.Sum([]byte(passcode + current.path))
		doors := [4]byte{digest[0] >> 4, digest[0] & 15, digest[1] >> 4, digest[1] & 15}
		directions := [4]grid.Position2D[int]{{Y: -1}, {Y: 1}, {X: -1}, {X: 1}}
		var next []vaultPath
		for i, direction := range directions {
			position := current.position.Add(direction)
			if doors[i] >= 11 && position.X >= 0 && position.Y >= 0 && position.X < 4 && position.Y < 4 {
				next = append(next, vaultPath{position, current.path + string("UDLR"[i])})
			}
		}
		return next
	}
}

func day17p01(r io.Reader) (string, error) {
	passcode, err := vaultPasscode(r)
	if err != nil {
		return "", err
	}
	for current := range search.BFSTree(vaultPath{}, vaultNeighbours(passcode)) {
		if current.atVault() {
			return current.path, nil
		}
	}
	return "", fmt.Errorf("vault is unreachable")
}

func day17p02(r io.Reader) (string, error) {
	passcode, err := vaultPasscode(r)
	if err != nil {
		return "", err
	}
	longest := -1
	for current := range search.DFSTree(vaultPath{}, vaultNeighbours(passcode)) {
		if current.atVault() {
			longest = max(longest, len(current.path))
		}
	}
	if longest < 0 {
		return "", fmt.Errorf("vault is unreachable")
	}
	return strconv.Itoa(longest), nil
}
