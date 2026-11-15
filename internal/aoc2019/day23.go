package aoc2019

import (
	"fmt"
	"io"
	"strconv"
)

const (
	networkSize    = 50
	networkNATAddr = 255
)

type packet struct {
	x, y int
}

type network struct {
	computers []*IntcodeComputer
	queues    [][]packet
}

func newNetwork(r io.Reader) (*network, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return nil, err
	}
	n := &network{
		computers: make([]*IntcodeComputer, networkSize),
		queues:    make([][]packet, networkSize),
	}
	for addr := range networkSize {
		n.computers[addr] = New(program)
		n.computers[addr].SetInput(addr)
	}
	return n, nil
}

// round lets every computer run until it waits for input once, and returns
// the packets sent to the NAT and whether any packet moved between computers.
func (n *network) round() ([]packet, bool, error) {
	var natPackets []packet
	active := false
	halted := 0
	for addr, computer := range n.computers {
		if computer.IsHalted() {
			halted++
			continue
		}
		if computer.IsWaiting() {
			if queue := n.queues[addr]; len(queue) > 0 {
				computer.AddInput(queue[0].x, queue[0].y)
				n.queues[addr] = queue[1:]
				active = true
			} else {
				computer.AddInput(-1)
			}
		}
		if err := computer.Run(); err != nil {
			return nil, false, fmt.Errorf("computer %d: %w", addr, err)
		}
		output := computer.GetOutput()
		if len(output)%3 != 0 {
			return nil, false, fmt.Errorf("computer %d sent a partial packet", addr)
		}
		for i := 0; i < len(output); i += 3 {
			dest, p := output[i], packet{x: output[i+1], y: output[i+2]}
			switch {
			case dest == networkNATAddr:
				natPackets = append(natPackets, p)
			case dest >= 0 && dest < networkSize:
				n.queues[dest] = append(n.queues[dest], p)
			default:
				return nil, false, fmt.Errorf("computer %d sent a packet to unknown address %d", addr, dest)
			}
			active = true
		}
	}
	if halted == networkSize {
		return nil, false, fmt.Errorf("every computer halted")
	}
	return natPackets, active, nil
}

func day23p01(r io.Reader) (string, error) {
	n, err := newNetwork(r)
	if err != nil {
		return "", err
	}
	for {
		natPackets, _, err := n.round()
		if err != nil {
			return "", err
		}
		if len(natPackets) > 0 {
			return strconv.Itoa(natPackets[0].y), nil
		}
	}
}

func day23p02(r io.Reader) (string, error) {
	n, err := newNetwork(r)
	if err != nil {
		return "", err
	}
	var nat packet
	hasNAT, delivered := false, false
	lastY := 0
	for {
		natPackets, active, err := n.round()
		if err != nil {
			return "", err
		}
		if len(natPackets) > 0 {
			nat, hasNAT = natPackets[len(natPackets)-1], true
		}
		if active || !hasNAT {
			continue
		}
		if delivered && nat.y == lastY {
			return strconv.Itoa(nat.y), nil
		}
		n.queues[0] = append(n.queues[0], nat)
		lastY, delivered = nat.y, true
	}
}
