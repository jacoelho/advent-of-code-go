package aoc2018

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"slices"
)

type trackPosition struct{ x, y int }
type mineCart struct {
	position        trackPosition
	direction, turn int
	alive           bool
}
type cartGoal uint8

const (
	firstCartCrash cartGoal = iota
	lastCartStanding
)

func mineCarts(r io.Reader, goal cartGoal) (string, error) {
	// Rows keep their own lengths because trailing spaces may be trimmed.
	var tracks [][]byte
	var carts []mineCart
	s := bufio.NewScanner(r)
	for s.Scan() {
		row := []byte(s.Text())
		for x, cell := range row {
			position := trackPosition{x, len(tracks)}
			direction := -1
			switch cell {
			case '^':
				direction, row[x] = 0, '|'
			case '>':
				direction, row[x] = 1, '-'
			case 'v':
				direction, row[x] = 2, '|'
			case '<':
				direction, row[x] = 3, '-'
			case ' ', '-', '|', '/', '\\', '+':
			default:
				return "", fmt.Errorf("invalid track cell %q", cell)
			}
			if direction >= 0 {
				carts = append(carts, mineCart{position: position, direction: direction, alive: true})
			}
		}
		tracks = append(tracks, row)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(carts) == 0 {
		return "", fmt.Errorf("track has no carts")
	}
	track := func(p trackPosition) byte {
		if p.y < 0 || p.y >= len(tracks) || p.x < 0 || p.x >= len(tracks[p.y]) {
			return ' '
		}
		return tracks[p.y][p.x]
	}

	occupied := make(map[trackPosition]int)
	for id, cart := range carts {
		occupied[cart.position] = id
	}
	directions := [...]trackPosition{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	order := make([]int, 0, len(carts))
	for {
		if len(occupied) == 0 {
			return "", fmt.Errorf("all mine carts collided")
		}
		if len(occupied) == 1 {
			if goal == firstCartCrash {
				return "", fmt.Errorf("fewer than two carts remain")
			}
			for position := range occupied {
				return fmt.Sprintf("%d,%d", position.x, position.y), nil
			}
		}

		order = order[:0]
		for id, cart := range carts {
			if cart.alive {
				order = append(order, id)
			}
		}
		slices.SortFunc(order, func(a, b int) int {
			return cmp.Or(cmp.Compare(carts[a].position.y, carts[b].position.y), cmp.Compare(carts[a].position.x, carts[b].position.x))
		})
		for _, id := range order {
			cart := &carts[id]
			if !cart.alive {
				continue
			}
			delete(occupied, cart.position)
			delta := directions[cart.direction]
			cart.position.x += delta.x
			cart.position.y += delta.y
			if other, collision := occupied[cart.position]; collision {
				cart.alive, carts[other].alive = false, false
				delete(occupied, cart.position)
				if goal == firstCartCrash {
					return fmt.Sprintf("%d,%d", cart.position.x, cart.position.y), nil
				}
				continue
			}
			switch track(cart.position) {
			case ' ':
				return "", fmt.Errorf("cart left the track at %d,%d", cart.position.x, cart.position.y)
			case '/':
				cart.direction ^= 1
			case '\\':
				cart.direction ^= 3
			case '+':
				cart.direction = (cart.direction + cart.turn + 3) % 4
				cart.turn = (cart.turn + 1) % 3
			case '-':
				if cart.direction%2 == 0 {
					return "", fmt.Errorf("vertical cart entered horizontal track")
				}
			case '|':
				if cart.direction%2 != 0 {
					return "", fmt.Errorf("horizontal cart entered vertical track")
				}
			}
			occupied[cart.position] = id
		}
	}
}

func day13p01(r io.Reader) (string, error) {
	return mineCarts(r, firstCartCrash)
}

func day13p02(r io.Reader) (string, error) {
	return mineCarts(r, lastCartStanding)
}
