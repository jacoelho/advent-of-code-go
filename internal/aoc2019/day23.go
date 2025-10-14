package aoc2019

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
)

type packet struct {
	x, y int
}

func initNetwork(program []int, numComputers int) ([]*IntcodeComputer, []*collections.Deque[packet]) {
	computers := make([]*IntcodeComputer, numComputers)
	queues := make([]*collections.Deque[packet], numComputers)

	for i := range numComputers {
		computers[i] = New(program)
		computers[i].SetInput(i)
		queues[i] = collections.NewDeque[packet](16)
	}

	return computers, queues
}

func processInput(computer *IntcodeComputer, queue *collections.Deque[packet]) bool {
	if p, ok := queue.PopFront(); ok {
		computer.AddInput(p.x, p.y)
		return true
	}
	computer.AddInput(-1)
	return false
}

func processOutput(
	computer *IntcodeComputer,
	queues []*collections.Deque[packet],
	numComputers int,
) (addr255Packets []packet, hadOutput bool) {
	output := computer.GetOutput()
	if len(output) < 3 {
		return nil, false
	}

	for outputPacket := range slices.Chunk(output, 3) {
		dest, x, y := outputPacket[0], outputPacket[1], outputPacket[2]

		if dest == 255 {
			addr255Packets = append(addr255Packets, packet{x: x, y: y})
		} else if dest >= 0 && dest < numComputers {
			queues[dest].PushBack(packet{x: x, y: y})
		}
	}
	return addr255Packets, true
}

func runNetwork(
	computers []*IntcodeComputer,
	queues []*collections.Deque[packet],
	numComputers int,
	shouldStop func(addr255Packets []packet, idle bool) (bool, int),
) (int, error) {
	for {
		idle := true
		var allAddr255Packets []packet

		for addr := range numComputers {
			computer := computers[addr]

			if computer.IsHalted() {
				continue
			}

			if computer.IsWaiting() {
				if processInput(computer, queues[addr]) {
					idle = false
				}
			}

			if err := computer.Run(); err != nil {
				return 0, err
			}

			addr255Packets, hadOutput := processOutput(computer, queues, numComputers)
			if hadOutput {
				idle = false
				allAddr255Packets = append(allAddr255Packets, addr255Packets...)
			}
		}

		if stop, result := shouldStop(allAddr255Packets, idle); stop {
			return result, nil
		}
	}
}

func day23p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	const numComputers = 50
	computers, queues := initNetwork(program, numComputers)

	result, err := runNetwork(computers, queues, numComputers, func(addr255Packets []packet, idle bool) (bool, int) {
		if len(addr255Packets) > 0 {
			return true, addr255Packets[0].y
		}
		return false, 0
	})
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}

func day23p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	const numComputers = 50
	computers, queues := initNetwork(program, numComputers)

	var natPacket *packet
	var lastSentY int
	natYSent := false

	result, err := runNetwork(computers, queues, numComputers, func(addr255Packets []packet, idle bool) (bool, int) {
		if len(addr255Packets) > 0 {
			natPacket = &addr255Packets[len(addr255Packets)-1]
		}

		if idle && natPacket != nil {
			if natYSent && natPacket.y == lastSentY {
				return true, natPacket.y
			}

			queues[0].PushBack(*natPacket)
			lastSentY = natPacket.y
			natYSent = true
		}

		return false, 0
	})
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}
