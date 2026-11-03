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

type vaultBudget struct{ paths, bytes int }

func (b *vaultBudget) consume(size int) error {
	if b.paths <= 0 || size > b.bytes {
		return fmt.Errorf("vault search exceeded its path or byte budget")
	}
	b.paths--
	b.bytes -= size
	return nil
}

func defaultVaultBudget() vaultBudget { return vaultBudget{1_000_000, 128 << 20} }

func vaultPasscode(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.ContainsAny(data, " \t\r\n") {
		return "", fmt.Errorf("invalid vault passcode")
	}
	return string(data), nil
}

func vaultNeighbours(passcode string, current vaultPath, budget *vaultBudget) ([]vaultPath, error) {
	if current.atVault() {
		return nil, nil
	}
	digest := md5.Sum([]byte(passcode + current.path))
	doors := [4]byte{digest[0] >> 4, digest[0] & 15, digest[1] >> 4, digest[1] & 15}
	directions := [4]grid.Position2D[int]{{Y: -1}, {Y: 1}, {X: -1}, {X: 1}}
	var next []vaultPath
	for i, direction := range directions {
		position := current.position.Add(direction)
		if doors[i] >= 11 && position.X >= 0 && position.Y >= 0 && position.X < 4 && position.Y < 4 {
			if err := budget.consume(len(passcode) + len(current.path) + 1); err != nil {
				return nil, err
			}
			next = append(next, vaultPath{position, current.path + string("UDLR"[i])})
		}
	}
	return next, nil
}

func day17p01(r io.Reader) (string, error) {
	passcode, err := vaultPasscode(r)
	if err != nil {
		return "", err
	}
	return shortestVault(passcode, defaultVaultBudget())
}

func shortestVault(passcode string, budget vaultBudget) (string, error) {
	if err := budget.consume(len(passcode)); err != nil {
		return "", err
	}
	children := func(current vaultPath) ([]vaultPath, error) {
		return vaultNeighbours(passcode, current, &budget)
	}
	for current, err := range search.BFSTree(vaultPath{}, children) {
		if err != nil {
			return "", err
		}
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
	longest, err := longestVault(passcode, defaultVaultBudget())
	if err != nil {
		return "", err
	}
	return strconv.Itoa(longest), nil
}

func longestVault(passcode string, budget vaultBudget) (int, error) {
	if err := budget.consume(len(passcode)); err != nil {
		return 0, err
	}
	children := func(current vaultPath) ([]vaultPath, error) {
		return vaultNeighbours(passcode, current, &budget)
	}
	longest := -1
	for current, err := range search.DFSTree(vaultPath{}, children) {
		if err != nil {
			return 0, err
		}
		if current.atVault() {
			longest = max(longest, len(current.path))
		}
	}
	if longest < 0 {
		return 0, fmt.Errorf("vault is unreachable")
	}
	return longest, nil
}
