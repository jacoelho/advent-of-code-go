package aoc2015

import (
	"cmp"
	"fmt"
	"io"
	"math/big"
	"math/bits"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func balancePackages(r io.Reader, groups int) (string, error) {
	s := scanner.NewScanner(r, convert.ScanNumber[int])
	weights := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(weights) < groups || len(weights) > 64 {
		return "", fmt.Errorf("need %d to 64 packages", groups)
	}
	total := 0
	for _, weight := range weights {
		if weight <= 0 || weight > int(^uint(0)>>1)-total {
			return "", fmt.Errorf("invalid package weight %d", weight)
		}
		total += weight
	}
	if total%groups != 0 {
		return "", fmt.Errorf("total weight is not divisible by %d", groups)
	}
	target := total / groups
	slices.SortFunc(weights, func(a, b int) int { return cmp.Compare(b, a) })
	if weights[0] > target {
		return "", fmt.Errorf("package exceeds group weight")
	}
	all := uint64(1)<<len(weights) - 1
	failed := make(map[uint64]bool)
	var partition func(uint64, int) bool
	partition = func(remaining uint64, count int) bool {
		if count == 1 {
			return true
		}
		if failed[remaining] {
			return false
		}
		first := bits.TrailingZeros64(remaining)
		// Groups are interchangeable; fix one package to avoid enumerating their permutations.
		var fill func(int, int, uint64) bool
		fill = func(start, needed int, chosen uint64) bool {
			if needed == 0 {
				return partition(remaining&^chosen, count-1)
			}
			for i := start; i < len(weights); i++ {
				if remaining&(1<<i) == 0 || weights[i] > needed {
					continue
				}
				if fill(i+1, needed-weights[i], chosen|1<<i) {
					return true
				}
			}
			return false
		}
		if fill(first+1, target-weights[first], 1<<first) {
			return true
		}
		failed[remaining] = true
		return false
	}
	for size := 1; size <= len(weights)/groups; size++ {
		var best *big.Int
		var choose func(int, int, int, uint64, *big.Int)
		choose = func(start, count, needed int, chosen uint64, product *big.Int) {
			if count == 0 {
				if needed == 0 && (best == nil || product.Cmp(best) < 0) && partition(all&^chosen, groups-1) {
					best = new(big.Int).Set(product)
				}
				return
			}
			if needed <= 0 || len(weights)-start < count {
				return
			}
			for i := start; i <= len(weights)-count; i++ {
				if weights[i] > needed {
					continue
				}
				next := new(big.Int).Mul(product, big.NewInt(int64(weights[i])))
				if best != nil && next.Cmp(best) >= 0 {
					continue
				}
				choose(i+1, count-1, needed-weights[i], chosen|1<<i, next)
			}
		}
		choose(0, size, target, 0, big.NewInt(1))
		if best != nil {
			return best.String(), nil
		}
	}
	return "", fmt.Errorf("packages cannot form %d equal groups", groups)
}

func day24p01(r io.Reader) (string, error) {
	return balancePackages(r, 3)
}

func day24p02(r io.Reader) (string, error) {
	return balancePackages(r, 4)
}
