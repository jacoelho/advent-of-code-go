package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type firewallLayer struct {
	depth int64
	width int64
}

func parseFirewall(r io.Reader) ([]firewallLayer, error) {
	s := scanner.NewScanner(r, func(line []byte) (firewallLayer, error) {
		depthText, widthText, ok := bytes.Cut(line, []byte(":"))
		depth, depthErr := convert.ScanNumber[int32](bytes.TrimSpace(depthText))
		width, widthErr := convert.ScanNumber[int32](bytes.TrimSpace(widthText))
		if !ok || depthErr != nil || widthErr != nil || depth < 0 || width < 1 {
			return firewallLayer{}, fmt.Errorf("invalid firewall layer %q", line)
		}
		return firewallLayer{int64(depth), int64(width)}, nil
	})

	var layers []firewallLayer
	seen := make(map[int64]bool)
	for layer := range s.Values() {
		if seen[layer.depth] {
			return nil, fmt.Errorf("duplicate firewall depth %d", layer.depth)
		}
		seen[layer.depth] = true
		layers = append(layers, layer)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(layers) == 0 {
		return nil, fmt.Errorf("empty firewall")
	}

	return layers, nil
}

func day13p01(r io.Reader) (string, error) {
	layers, err := parseFirewall(r)
	if err != nil {
		return "", err
	}

	severity := new(big.Int)
	for _, layer := range layers {
		if layer.width == 1 || layer.depth%(2*(layer.width-1)) == 0 {
			severity.Add(severity, big.NewInt(layer.depth*layer.width))
		}
	}

	return severity.String(), nil
}

func day13p02(r io.Reader) (string, error) {
	layers, err := parseFirewall(r)
	if err != nil {
		return "", err
	}

	const searchLimit int64 = 100_000_000
	cycle := int64(1)
	completeCycle := true
	for _, layer := range layers {
		if layer.width == 1 {
			return "", fmt.Errorf("firewall layer %d always catches the packet", layer.depth)
		}
		if !completeCycle {
			continue
		}
		period := 2 * (layer.width - 1)
		factor := period / xmath.GCD(cycle, period)
		if cycle > searchLimit/factor {
			completeCycle = false
			cycle = searchLimit
			continue
		}
		cycle *= factor
	}

	// Each scanner rejects one residue class; mark delays in bounded blocks.
	var blocked [65536]bool
	for start := int64(0); start < cycle; start += int64(len(blocked)) {
		clear(blocked[:])
		end := min(start+int64(len(blocked)), cycle)
		for _, layer := range layers {
			period := 2 * (layer.width - 1)
			first := (period - layer.depth%period) % period
			if first < start {
				first += (start - first + period - 1) / period * period
			}
			for delay := first; delay < end; delay += period {
				blocked[delay-start] = true
			}
		}
		for offset := range end - start {
			if !blocked[offset] {
				return strconv.FormatInt(start+offset, 10), nil
			}
		}
	}

	if completeCycle {
		return "", fmt.Errorf("firewall has no safe delay")
	}
	return "", fmt.Errorf("firewall search exceeds %d delays", searchLimit)
}
