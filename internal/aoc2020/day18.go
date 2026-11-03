package aoc2020

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type homeworkExpression struct {
	text             string
	position         int
	additionPriority int
}

func (e *homeworkExpression) skipSpaces() {
	for e.position < len(e.text) && (e.text[e.position] == ' ' || e.text[e.position] == '\t') {
		e.position++
	}
}

func (e *homeworkExpression) operand() (int64, error) {
	e.skipSpaces()
	if e.position == len(e.text) {
		return 0, fmt.Errorf("missing operand at column %d", e.position+1)
	}
	if e.text[e.position] == '(' {
		e.position++
		value, err := e.expression(1)
		if err != nil {
			return 0, err
		}
		e.skipSpaces()
		if e.position == len(e.text) || e.text[e.position] != ')' {
			return 0, fmt.Errorf("missing closing parenthesis at column %d", e.position+1)
		}
		e.position++
		return value, nil
	}
	start := e.position
	for e.position < len(e.text) && e.text[e.position] >= '0' && e.text[e.position] <= '9' {
		e.position++
	}
	if start == e.position {
		return 0, fmt.Errorf("expected number at column %d", start+1)
	}
	value, err := strconv.ParseInt(e.text[start:e.position], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("number at column %d: %w", start+1, err)
	}
	return value, nil
}

func (e *homeworkExpression) expression(minimumPriority int) (int64, error) {
	value, err := e.operand()
	if err != nil {
		return 0, err
	}
	for {
		e.skipSpaces()
		if e.position == len(e.text) || e.text[e.position] == ')' {
			return value, nil
		}
		operator := e.text[e.position]
		priority := 1
		switch operator {
		case '+':
			priority = e.additionPriority
		case '*':
		default:
			return 0, fmt.Errorf("invalid operator %q at column %d", operator, e.position+1)
		}
		if priority < minimumPriority {
			return value, nil
		}
		e.position++
		// Raising the right operand's minimum makes operators of equal precedence left associative.
		right, err := e.expression(priority + 1)
		if err != nil {
			return 0, err
		}
		if operator == '+' {
			if value > math.MaxInt64-right {
				return 0, fmt.Errorf("addition overflows int64")
			}
			value += right
		} else {
			if right != 0 && value > math.MaxInt64/right {
				return 0, fmt.Errorf("multiplication overflows int64")
			}
			value *= right
		}
	}
}

func sumHomework(r io.Reader, additionPriority int) (string, error) {
	s := scanner.NewScanner(r, func(line []byte) (int64, error) {
		expression := homeworkExpression{text: string(line), additionPriority: additionPriority}
		value, err := expression.expression(1)
		if err != nil {
			return 0, err
		}
		if expression.position != len(expression.text) {
			return 0, fmt.Errorf("unexpected token at column %d", expression.position+1)
		}
		return value, nil
	})
	var total int64
	for value := range s.Values() {
		if total > math.MaxInt64-value {
			return "", fmt.Errorf("homework total overflows int64")
		}
		total += value
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}

func day18p01(r io.Reader) (string, error) {
	return sumHomework(r, 1)
}

func day18p02(r io.Reader) (string, error) {
	return sumHomework(r, 2)
}
