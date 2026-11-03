package aoc2022

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type monkeyJob struct {
	number      *big.Int
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
		if _, exists := jobs[name]; exists || len(jobs) == 10_000 {
			return nil, nil, fmt.Errorf("duplicate monkey or more than 10000 definitions")
		}
		fields := strings.Fields(text)
		var job monkeyJob
		switch len(fields) {
		case 1:
			number, err := convert.ScanNumber[int64]([]byte(fields[0]))
			if err != nil {
				return nil, nil, fmt.Errorf("monkey %s: %w", name, err)
			}
			job.number = big.NewInt(number)
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

	marks, heights := make(map[string]byte), make(map[string]int)
	var order []string
	var visit func(string, int) error
	visit = func(name string, depth int) error {
		if depth > 1024 {
			return fmt.Errorf("monkey expression exceeds 1024 levels")
		}
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
		marks[name], heights[name] = 1, 1
		if job.number == nil {
			if err := visit(job.left, depth+1); err != nil {
				return err
			}
			if err := visit(job.right, depth+1); err != nil {
				return err
			}
			heights[name] += max(heights[job.left], heights[job.right])
			if heights[name] > 1024 {
				return fmt.Errorf("monkey expression exceeds 1024 levels")
			}
		}
		marks[name] = 2
		order = append(order, name)
		return nil
	}
	for _, name := range names {
		if err := visit(name, 1); err != nil {
			return nil, nil, err
		}
	}
	return jobs, order, nil
}

func monkeyArithmetic(operation byte, left, right *big.Int) (*big.Int, error) {
	result := new(big.Int)
	switch operation {
	case '+':
		result.Add(left, right)
	case '-':
		result.Sub(left, right)
	case '*':
		result.Mul(left, right)
	case '/':
		if right.Sign() == 0 {
			return nil, fmt.Errorf("monkey division by zero")
		}
		remainder := new(big.Int)
		result.QuoRem(left, right, remainder)
		if remainder.Sign() != 0 {
			return nil, fmt.Errorf("monkey division is not exact")
		}
	}
	if result.BitLen() > 4096 {
		return nil, fmt.Errorf("monkey arithmetic exceeds 4096 bits")
	}
	return result, nil
}

type monkeyEvaluation struct {
	jobs   map[string]monkeyJob
	values map[string]*big.Int
	human  *big.Int
}

func (e *monkeyEvaluation) value(name string) (*big.Int, error) {
	if name == "humn" && e.human != nil {
		return e.human, nil
	}
	if value, exists := e.values[name]; exists {
		return value, nil
	}
	job := e.jobs[name]
	if job.number != nil {
		return job.number, nil
	}
	left, err := e.value(job.left)
	if err != nil {
		return nil, err
	}
	right, err := e.value(job.right)
	if err != nil {
		return nil, err
	}
	value, err := monkeyArithmetic(job.operation, left, right)
	if err != nil {
		return nil, fmt.Errorf("monkey %s: %w", name, err)
	}
	e.values[name] = value
	return value, nil
}

func day21p01(r io.Reader) (string, error) {
	jobs, _, err := monkeyJobs(r)
	if err != nil {
		return "", err
	}
	evaluation := monkeyEvaluation{jobs: jobs, values: make(map[string]*big.Int)}
	value, err := evaluation.value("root")
	if err != nil {
		return "", err
	}
	return value.String(), nil
}

func day21p02(r io.Reader) (string, error) {
	jobs, order, err := monkeyJobs(r)
	if err != nil {
		return "", err
	}
	human, exists := jobs["humn"]
	if !exists || human.number == nil || jobs["root"].number != nil {
		return "", fmt.Errorf("equation requires a human number and a root operation")
	}
	dependent := make(map[string]bool)
	for _, name := range order {
		job := jobs[name]
		dependent[name] = name == "humn" || job.number == nil && (dependent[job.left] || dependent[job.right])
	}
	root := jobs["root"]
	if dependent[root.left] == dependent[root.right] {
		return "", fmt.Errorf("root equality requires exactly one human-dependent side")
	}
	name, known := root.left, root.right
	if dependent[root.right] {
		name, known = root.right, root.left
	}
	evaluation := monkeyEvaluation{jobs: jobs, values: make(map[string]*big.Int)}
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
			if constant.Sign() == 0 {
				if target.Sign() == 0 {
					return "", fmt.Errorf("unsupported noninvertible monkey equation: multiplication by zero")
				}
				return "", fmt.Errorf("monkey multiplication by zero has no solution")
			}
			target, err = monkeyArithmetic('/', target, constant)
		case '/':
			if left {
				if constant.Sign() == 0 {
					return "", fmt.Errorf("monkey division by zero")
				}
				target, err = monkeyArithmetic('*', target, constant)
			} else {
				if target.Sign() == 0 {
					if constant.Sign() == 0 {
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

	check := monkeyEvaluation{jobs: jobs, values: make(map[string]*big.Int), human: target}
	left, err := check.value(root.left)
	if err != nil {
		return "", err
	}
	right, err := check.value(root.right)
	if err != nil {
		return "", err
	}
	if left.Cmp(right) != 0 {
		return "", fmt.Errorf("human candidate does not satisfy root equality")
	}
	return target.String(), nil
}
