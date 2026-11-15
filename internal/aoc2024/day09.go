package aoc2024

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func parseDiskMap(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSuffix(data, []byte("\n"))
	data = bytes.TrimSuffix(data, []byte("\r"))
	if len(data) == 0 {
		return nil, fmt.Errorf("disk map is empty")
	}
	for _, digit := range data {
		if digit < '0' || digit > '9' {
			return nil, fmt.Errorf("invalid disk map digit %q", digit)
		}
	}
	return data, nil
}

func day09p01(r io.Reader) (string, error) {
	digits, err := parseDiskMap(r)
	if err != nil {
		return "", err
	}
	var disk []int
	for i, digit := range digits {
		value := -1
		if i%2 == 0 {
			value = i / 2
		}
		for range int(digit - '0') {
			disk = append(disk, value)
		}
	}

	front, back := 0, len(disk)-1
	for front < back {
		if disk[front] == -1 && disk[back] != -1 {
			disk[front], disk[back] = disk[back], disk[front]
		}
		if disk[front] != -1 {
			front++
		}
		if disk[back] == -1 {
			back--
		}
	}

	var total int64
	for position, fileID := range disk {
		if fileID != -1 {
			total += int64(position) * int64(fileID)
		}
	}
	return strconv.FormatInt(total, 10), nil
}

type diskSpan struct {
	index int
	size  int
}

func diskSpans(digits []byte) ([]diskSpan, []diskSpan) {
	var files, free []diskSpan
	position := 0
	for i, digit := range digits {
		size := int(digit - '0')
		if i%2 == 0 {
			files = append(files, diskSpan{index: position, size: size})
		} else if size != 0 {
			last := len(free) - 1
			if last >= 0 && free[last].index+free[last].size == position {
				free[last].size += size
			} else {
				free = append(free, diskSpan{index: position, size: size})
			}
		}
		position += size
	}
	return files, free
}

func day09p02(r io.Reader) (string, error) {
	digits, err := parseDiskMap(r)
	if err != nil {
		return "", err
	}
	files, free := diskSpans(digits)
	// Free spans only shrink from the left and files move right to left, so a
	// span too small for a size, or not left of the current file, stays unusable
	// for that size: the first fit for each size never moves left.
	var firstFit [10]int
	for fileID := len(files) - 1; fileID > 0; fileID-- {
		file := files[fileID]
		if file.size == 0 {
			continue
		}
		i := firstFit[file.size]
		for i < len(free) && free[i].index < file.index && free[i].size < file.size {
			i++
		}
		firstFit[file.size] = i
		if i == len(free) || free[i].index >= file.index {
			continue
		}
		// The vacated span lies right of every file still to be moved, so it is
		// never added to the free list.
		files[fileID].index = free[i].index
		free[i].index += file.size
		free[i].size -= file.size
	}

	var total int64
	for fileID, file := range files {
		for offset := range file.size {
			total += int64(fileID) * int64(file.index+offset)
		}
	}
	return strconv.FormatInt(total, 10), nil
}
