package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type firewallLayer struct {
	depth int64
	width int64
}

// period is how many picoseconds the layer's scanner takes to return to the top.
func (layer firewallLayer) period() int64 {
	return 2 * (layer.width - 1)
}

func parseFirewall(r io.Reader) ([]firewallLayer, error) {
	s := scanner.NewScanner(r, func(line []byte) (firewallLayer, error) {
		depthText, widthText, ok := bytes.Cut(line, []byte(":"))
		depth, depthErr := convert.ScanNumber[int64](bytes.TrimSpace(depthText))
		width, widthErr := convert.ScanNumber[int64](bytes.TrimSpace(widthText))
		if !ok || depthErr != nil || widthErr != nil || depth < 0 || width < 1 {
			return firewallLayer{}, fmt.Errorf("invalid firewall layer %q", line)
		}
		if _, ok := xmath.CheckedMul(width-1, 2); !ok {
			return firewallLayer{}, fmt.Errorf("firewall layer %q has a period beyond int64", line)
		}
		return firewallLayer{depth, width}, nil
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

	var severity int64
	for _, layer := range layers {
		if layer.width > 1 && layer.depth%layer.period() != 0 {
			continue
		}
		damage, damageOK := xmath.CheckedMul(layer.depth, layer.width)
		total, totalOK := xmath.CheckedAdd(severity, damage)
		if !damageOK || !totalOK {
			return "", fmt.Errorf("firewall severity overflows int64")
		}
		severity = total
	}

	return strconv.FormatInt(severity, 10), nil
}

func day13p02(r io.Reader) (string, error) {
	layers, err := parseFirewall(r)
	if err != nil {
		return "", err
	}

	// Delays repeat with the LCM of the periods; past int64 the search is unbounded in practice.
	cycle := int64(1)
	for _, layer := range layers {
		if layer.width == 1 {
			return "", fmt.Errorf("firewall layer %d always catches the packet", layer.depth)
		}
		factor := layer.period() / xmath.GCD(cycle, layer.period())
		next, ok := xmath.CheckedMul(cycle, factor)
		if !ok {
			next = math.MaxInt64
		}
		cycle = next
	}

	// Each scanner rejects one residue class; mark delays in bounded blocks.
	var blocked [65536]bool
	for start := int64(0); start < cycle; start += int64(len(blocked)) {
		clear(blocked[:])
		size := min(int64(len(blocked)), cycle-start)
		for _, layer := range layers {
			period := layer.period()
			caught := (period - layer.depth%period) % period
			for offset := xmath.Modulo(caught-start, period); offset < size; offset += min(period, size) {
				blocked[offset] = true
			}
		}
		for offset := range size {
			if !blocked[offset] {
				return strconv.FormatInt(start+offset, 10), nil
			}
		}
		if size < int64(len(blocked)) {
			break
		}
	}

	return "", fmt.Errorf("firewall has no safe delay")
}
