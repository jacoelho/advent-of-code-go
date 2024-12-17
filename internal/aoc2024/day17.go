package aoc2024

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

func parseComputerProgram(r io.Reader) (*computer, error) {
	s := bufio.NewScanner(r)

	var digits []int
	for s.Scan() {
		digits = append(digits, convert.ExtractDigits[int](s.Text())...)
	}
	if s.Err() != nil || len(digits) == 0 {
		return nil, errors.New("invalid input")
	}
	return &computer{
		a:            digits[0],
		b:            digits[1],
		c:            digits[2],
		instructions: digits[3:],
	}, nil
}

type computer struct {
	a            int
	b            int
	c            int
	instructions []int
	ip           int
}

func (c *computer) clone() *computer {
	return &computer{
		a:            c.a,
		b:            c.b,
		c:            c.c,
		instructions: slices.Clone(c.instructions),
	}
}

func (c *computer) combo(v int) int {
	switch v {
	case 0, 1, 2, 3:
		return v
	case 4:
		return c.a
	case 5:
		return c.b
	case 6:
		return c.c
	case 7:
		panic("reserved")
	default:
		panic("unreachable")
	}
}

func (c *computer) run() iter.Seq[int] {
	return func(yield func(int) bool) {
		for c.ip < len(c.instructions) {
			instruction := c.instructions[c.ip]
			operand := c.instructions[c.ip+1]

			switch instruction {
			case 0:
				c.a = c.a >> c.combo(operand)
			case 1:
				c.b = c.b ^ operand
			case 2:
				// a % b is equivalent to (b - 1) & a if b is power of 2
				c.b = 7 & c.combo(operand)
			case 3:
				if c.a != 0 {
					c.ip = operand
					continue
				}
			case 4:
				c.b = c.b ^ c.c
			case 5:
				if !yield(7 & c.combo(operand)) {
					return
				}
			case 6:
				c.b = c.a >> c.combo(operand)
			case 7:
				c.c = c.a >> c.combo(operand)
			}
			c.ip += 2
		}
	}
}

//lint:file-ignore U1000 debug function
func (c *computer) debug() string {
	sb := new(strings.Builder)

	combo := func(v int) string {
		switch v {
		case 4:
			return "A"
		case 5:
			return "B"
		case 6:
			return "C"
		default:
			return strconv.Itoa(v)
		}
	}

	instructionMap := map[int]string{
		0: "A = A >> %s",
		1: "B = B ^ %s",
		2: "B = %s %% 8",
		3: "IF A != 0 JMP %d",
		4: "B = B ^ C",
		5: "OUTPUT %s %% 8",
		6: "B = A >> %s",
		7: "C = A >> %s",
	}

	for ip, v := range xiter.Enumerate(slices.Chunk(c.instructions, 2)) {
		instruction := v[0]
		operand := v[1]

		format := instructionMap[instruction]
		switch instruction {
		case 3:
			fmt.Fprintf(sb, "%d: "+format+"\n", ip, operand)
		case 4:
			fmt.Fprintf(sb, "%d: "+format+"\n", ip)
		default:
			fmt.Fprintf(sb, "%d: "+format+"\n", ip, combo(operand))
		}
	}

	return sb.String()
}

func day17p01(r io.Reader) (string, error) {
	c := aoc.Must(parseComputerProgram(r))

	s := xiter.Reduce(func(sum []string, e int) []string {
		return append(sum, strconv.Itoa(e))
	}, []string{}, c.run())

	return strings.Join(s, ","), nil
}

func day17p02(r io.Reader) (string, error) {
	c := aoc.Must(parseComputerProgram(r))
	i := 1
	for {
		cc := c.clone()
		cc.a = i
		result := slices.Collect(cc.run())

		switch {
		case slices.Equal(result, c.instructions):
			return strconv.Itoa(i), nil

		// digit matches
		// attempt next
		case slices.Equal(result, c.instructions[len(c.instructions)-len(result):]):
			i <<= 3 // i *=8

		default:
			i++
		}
	}
}
