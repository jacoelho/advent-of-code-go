package aoc2022

import (
	"cmp"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type distressPacket struct {
	number int64
	list   []distressPacket
	isList bool
}

func parsePackets(r io.Reader) ([]distressPacket, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return nil, err
	}
	var packets []distressPacket
	nodes := 0
	for _, block := range strings.Split(input, "\n\n") {
		lines := strings.Split(block, "\n")
		if len(lines) != 2 {
			return nil, fmt.Errorf("expected packet pair")
		}
		for _, line := range lines {
			at := 0
			var parse func(int) (distressPacket, error)
			parse = func(depth int) (distressPacket, error) {
				nodes++
				if depth > 128 || nodes > 1_000_000 || at >= len(line) {
					return distressPacket{}, fmt.Errorf("invalid or oversized packet")
				}
				if line[at] != '[' {
					start := at
					for at < len(line) && line[at] >= '0' && line[at] <= '9' {
						at++
					}
					value, err := convert.ScanNumber[int64]([]byte(line[start:at]))
					return distressPacket{number: value}, err
				}
				result := distressPacket{isList: true}
				at++
				if at < len(line) && line[at] == ']' {
					at++
					return result, nil
				}
				for {
					child, err := parse(depth + 1)
					if err != nil {
						return distressPacket{}, err
					}
					result.list = append(result.list, child)
					if at >= len(line) {
						return distressPacket{}, fmt.Errorf("unclosed packet")
					}
					switch line[at] {
					case ']':
						at++
						return result, nil
					case ',':
						at++
					default:
						return distressPacket{}, fmt.Errorf("invalid packet delimiter")
					}
				}
			}
			packet, err := parse(0)
			if err != nil {
				return nil, err
			}
			if at != len(line) || !packet.isList {
				return nil, fmt.Errorf("expected complete list packet")
			}
			packets = append(packets, packet)
		}
	}
	return packets, nil
}

func comparePackets(a, b distressPacket) int {
	if !a.isList && !b.isList {
		return cmp.Compare(a.number, b.number)
	}
	left, right := a.list, b.list
	if !a.isList {
		left = []distressPacket{a}
	}
	if !b.isList {
		right = []distressPacket{b}
	}
	for i := range min(len(left), len(right)) {
		if order := comparePackets(left[i], right[i]); order != 0 {
			return order
		}
	}
	return cmp.Compare(len(left), len(right))
}

func day13p01(r io.Reader) (string, error) {
	packets, err := parsePackets(r)
	if err != nil {
		return "", err
	}
	total := 0
	for i := 0; i < len(packets); i += 2 {
		if comparePackets(packets[i], packets[i+1]) < 0 {
			total += i/2 + 1
		}
	}
	return strconv.Itoa(total), nil
}

func day13p02(r io.Reader) (string, error) {
	packets, err := parsePackets(r)
	if err != nil {
		return "", err
	}
	first, second := 1, 2
	// Inserted dividers follow existing equal packets in stable order.
	for _, packet := range packets {
		if comparePackets(packet, distressPacket{number: 2}) <= 0 {
			first++
		}
		if comparePackets(packet, distressPacket{number: 6}) <= 0 {
			second++
		}
	}
	return strconv.Itoa(first * second), nil
}
