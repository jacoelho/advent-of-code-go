package aoc2024

import (
	"bufio"
	"io"
	"maps"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/xiter"
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

type block struct {
	index int
	size  int
}

func parseDiskAsBlocks(r io.Reader) (map[int]block, []block) {
	var (
		id     int
		index  int
		blanks []block
	)
	result := make(map[int]block)
	s := bufio.NewScanner(r)
	for s.Scan() {
		for i, v := range s.Bytes() {
			size := int(v - '0')
			if i%2 == 0 {
				result[id] = block{index, size}
				id++
			} else {
				blanks = append(blanks, block{index, size})
			}
			index += size
		}
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
	return result, blanks
}

func day09p02(r io.Reader) (string, error) {
	blocks, blanks := parseDiskAsBlocks(r)

	for fileID := xiter.Max(maps.Keys(blocks)); fileID > 0; fileID-- {
		fileBlock := blocks[fileID]

		for i, blank := range blanks {
			if blank.index >= fileBlock.index {
				blanks = blanks[:i]
				break
			}
			if blank.size < fileBlock.size {
				continue
			}

			blocks[fileID] = block{index: blank.index, size: fileBlock.size}

			if blank.size == fileBlock.size {
				blanks[i] = block{index: blank.index, size: 0}

				// we could also delete, but is it slower
				// slices.Delete(blanks, i, i+1)
			} else {
				blanks[i] = block{
					index: blank.index + fileBlock.size,
					size:  blank.size - fileBlock.size,
				}
			}
			break
		}
	}

	var total int
	for fileID, fileBlock := range blocks {
		start := fileBlock.index
		end := fileBlock.index + fileBlock.size - 1
		total += fileID * (end*(end+1) - start*(start-1)) / 2
	}

	return strconv.Itoa(total), nil
}
