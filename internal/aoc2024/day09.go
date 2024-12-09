package aoc2024

import (
	"bufio"
	"io"
	"slices"
	"strconv"
)

func parseDiskAsSparse(r io.Reader) []int {
	var (
		id     int
		result []int
	)
	s := bufio.NewScanner(r)
	for s.Scan() {
		for i, v := range s.Bytes() {
			value := -1
			if i%2 == 0 {
				value = id
				id++
			}
			result = append(result, slices.Repeat([]int{value}, int(v-'0'))...)
		}
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
	return result
}

func day09p01(r io.Reader) (string, error) {
	disk := parseDiskAsSparse(r)

	frontPtr := 0
	backPtr := len(disk) - 1
	for frontPtr < backPtr {
		if disk[frontPtr] == -1 && disk[backPtr] != -1 {
			disk[frontPtr], disk[backPtr] = disk[backPtr], disk[frontPtr]
			frontPtr++
			backPtr--
		}
		if disk[frontPtr] != -1 {
			frontPtr++
		}
		if disk[backPtr] == -1 {
			backPtr--
		}
	}

	var total int
	for i, v := range disk {
		if v != -1 {
			total += i * v
		}
	}
	return strconv.Itoa(total), nil
}
