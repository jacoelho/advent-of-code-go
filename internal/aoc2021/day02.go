package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type submarineNavigation byte

const (
	directNavigation submarineNavigation = iota
	aimedNavigation
)

func navigationSum(a, b int64) (int64, error) {
	if b > 0 && a > math.MaxInt64-b || b < 0 && a < math.MinInt64-b {
		return 0, fmt.Errorf("navigation addition overflows int64")
	}
	return a + b, nil
}

func navigationProduct(a, b int64) (int64, error) {
	product := a * b
	if a == -1 && b == math.MinInt64 || b == -1 && a == math.MinInt64 || a != 0 && product/a != b {
		return 0, fmt.Errorf("navigation multiplication overflows int64")
	}
	return product, nil
}

func submarineCourse(r io.Reader, mode submarineNavigation) (int64, error) {
	s := bufio.NewScanner(r)
	var horizontal, depth, aim int64
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) != 2 {
			return 0, fmt.Errorf("invalid submarine command %q", s.Text())
		}
		distance, err := convert.ScanNumber[int64]([]byte(fields[1]))
		if err != nil {
			return 0, err
		}
		if distance < 0 {
			return 0, fmt.Errorf("negative command distance")
		}
		switch fields[0] {
		case "forward":
			horizontal, err = navigationSum(horizontal, distance)
			if err == nil && mode == aimedNavigation {
				var change int64
				change, err = navigationProduct(aim, distance)
				if err == nil {
					depth, err = navigationSum(depth, change)
				}
			}
		case "up", "down":
			if fields[0] == "up" {
				distance = -distance
			}
			if mode == directNavigation {
				depth, err = navigationSum(depth, distance)
			} else {
				aim, err = navigationSum(aim, distance)
			}
		default:
			return 0, fmt.Errorf("unknown submarine command %q", fields[0])
		}
		if err != nil {
			return 0, err
		}
	}
	if err := s.Err(); err != nil {
		return 0, err
	}
	return navigationProduct(horizontal, depth)
}

func day02p01(r io.Reader) (string, error) {
	result, err := submarineCourse(r, directNavigation)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(result, 10), nil
}

func day02p02(r io.Reader) (string, error) {
	result, err := submarineCourse(r, aimedNavigation)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(result, 10), nil
}
