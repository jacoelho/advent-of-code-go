package aoc2022

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

// monkeyJob is a number job when operation is zero.
type monkeyJob struct {
	number      int64
	left, right string
	operation   byte
}

func monkeyJobs(r io.Reader) (map[string]monkeyJob, []string, error) {
	s := bufio.NewScanner(r)
	namePattern := regexp.MustCompile(`^[a-z]{4}$`)
	jobs := make(map[string]monkeyJob)
	var names []string
	for s.Scan() {
		name, text, found := strings.Cut(s.Text(), ": ")
		if !found || !namePattern.MatchString(name) {
			return nil, nil, fmt.Errorf("invalid monkey definition %q", s.Text())
		}
		if _, exists := jobs[name]; exists {
			return nil, nil, fmt.Errorf("duplicate monkey %s", name)
		}
		fields := strings.Fields(text)
		var job monkeyJob
		switch len(fields) {
		case 1:
			number, err := convert.ScanNumber[int64]([]byte(fields[0]))
			if err != nil {
				return nil, nil, fmt.Errorf("monkey %s: %w", name, err)
			}
			job.number = number
		case 3:
			if !namePattern.MatchString(fields[0]) || !namePattern.MatchString(fields[2]) || len(fields[1]) != 1 || !strings.ContainsAny(fields[1], "+-*/") {
				return nil, nil, fmt.Errorf("invalid monkey operation %q", text)
			}
			job.left, job.right, job.operation = fields[0], fields[2], fields[1][0]
		default:
			return nil, nil, fmt.Errorf("invalid monkey job %q", text)
		}
		jobs[name] = job
		names = append(names, name)
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	if _, exists := jobs["root"]; !exists {
		return nil, nil, fmt.Errorf("missing root monkey")
	}

	marks := make(map[string]byte)
	var order []string
	var visit func(string) error
	visit = func(name string) error {
		if marks[name] == 1 {
			return fmt.Errorf("cycle at monkey %s", name)
		}
		if marks[name] == 2 {
			return nil
		}
		job, exists := jobs[name]
		if !exists {
			return fmt.Errorf("missing referenced monkey %s", name)
		}
		marks[name] = 1
		if job.operation != 0 {
			if err := visit(job.left); err != nil {
				return err
			}
			if err := visit(job.right); err != nil {
				return err
			}
		}
		marks[name] = 2
		order = append(order, name)
		return nil
	}
	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, nil, err
		}
	}
	return jobs, order, nil
}

func monkeyArithmetic(operation byte, left, right int64) (int64, error) {
	var result int64
	ok := true
	switch operation {
	case '+':
		result, ok = xmath.CheckedAdd(left, right)
	case '-':
		result, ok = xmath.CheckedSub(left, right)
	case '*':
		result, ok = xmath.CheckedMul(left, right)
	case '/':
		if right == 0 {
			return 0, fmt.Errorf("monkey division by zero")
		}
		if left == math.MinInt64 && right == -1 {
			ok = false
			break
		}
		if left%right != 0 {
			return 0, fmt.Errorf("monkey division %d / %d is not exact", left, right)
		}
		result = left / right
	}
	if !ok {
		return 0, fmt.Errorf("monkey arithmetic %d %c %d overflows int64", left, operation, right)
	}
	return result, nil
}

type monkeyEvaluation struct {
	jobs   map[string]monkeyJob
	values map[string]int64
}

func (e *monkeyEvaluation) value(name string) (int64, error) {
	if value, exists := e.values[name]; exists {
		return value, nil
	}
	job := e.jobs[name]
	if job.operation == 0 {
		return job.number, nil
	}
	left, err := e.value(job.left)
	if err != nil {
		return 0, err
	}
	right, err := e.value(job.right)
	if err != nil {
		return 0, err
	}
	value, err := monkeyArithmetic(job.operation, left, right)
	if err != nil {
		return 0, fmt.Errorf("monkey %s: %w", name, err)
	}
	e.values[name] = value
	return value, nil
}

func day21p01(r io.Reader) (string, error) {
	jobs, _, err := monkeyJobs(r)
	if err != nil {
		return "", err
	}
	evaluation := monkeyEvaluation{jobs: jobs, values: make(map[string]int64)}
	value, err := evaluation.value("root")
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(value, 10), nil
}

func day21p02(r io.Reader) (string, error) {
	jobs, order, err := monkeyJobs(r)
	if err != nil {
		return "", err
	}
	human, exists := jobs["humn"]
	if !exists || human.operation != 0 || jobs["root"].operation == 0 {
		return "", fmt.Errorf("equation requires a human number and a root operation")
	}
	dependent := make(map[string]bool)
	for _, name := range order {
		job := jobs[name]
		dependent[name] = name == "humn" || job.operation != 0 && (dependent[job.left] || dependent[job.right])
	}
	root := jobs["root"]
	if dependent[root.left] == dependent[root.right] {
		return "", fmt.Errorf("root equality requires exactly one human-dependent side")
	}
	name, known := root.left, root.right
	if dependent[root.right] {
		name, known = root.right, root.left
	}
	evaluation := monkeyEvaluation{jobs: jobs, values: make(map[string]int64)}
	target, err := evaluation.value(known)
	if err != nil {
		return "", err
	}
	for name != "humn" {
		job := jobs[name]
		left := dependent[job.left]
		if left == dependent[job.right] {
			return "", fmt.Errorf("monkey %s requires exactly one human-dependent operand", name)
		}
		next, other := job.left, job.right
		if !left {
			next, other = job.right, job.left
		}
		constant, err := evaluation.value(other)
		if err != nil {
			return "", err
		}
		switch job.operation {
		case '+':
			target, err = monkeyArithmetic('-', target, constant)
		case '-':
			if left {
				target, err = monkeyArithmetic('+', target, constant)
			} else {
				target, err = monkeyArithmetic('-', constant, target)
			}
		case '*':
			if constant == 0 {
				if target == 0 {
					return "", fmt.Errorf("unsupported noninvertible monkey equation: multiplication by zero")
				}
				return "", fmt.Errorf("monkey multiplication by zero has no solution")
			}
			target, err = monkeyArithmetic('/', target, constant)
		case '/':
			if left {
				if constant == 0 {
					return "", fmt.Errorf("monkey division by zero")
				}
				target, err = monkeyArithmetic('*', target, constant)
			} else {
				if target == 0 {
					if constant == 0 {
						return "", fmt.Errorf("unsupported noninvertible monkey equation: zero numerator and zero quotient")
					}
					return "", fmt.Errorf("nonzero numerator cannot give an exact zero quotient")
				}
				target, err = monkeyArithmetic('/', constant, target)
			}
		}
		if err != nil {
			return "", fmt.Errorf("invert monkey %s: %w", name, err)
		}
		name = next
	}

	check := monkeyEvaluation{jobs: jobs, values: map[string]int64{"humn": target}}
	left, err := check.value(root.left)
	if err != nil {
		return "", err
	}
	right, err := check.value(root.right)
	if err != nil {
		return "", err
	}
	if left != right {
		return "", fmt.Errorf("human candidate does not satisfy root equality")
	}
	return strconv.FormatInt(target, 10), nil
}
