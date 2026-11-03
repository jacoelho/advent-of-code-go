package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	internalDir = "internal"
	maxStdout   = 16 << 20
	maxStderr   = 1 << 20
)

var (
	yearDirRegex   = regexp.MustCompile(`^aoc(\d{4})$`)
	testNameRegex  = regexp.MustCompile(`^Test_day(\d{2})p(\d{2})$`)
	errOutputLimit = errors.New("test command output limit exceeded")
)

type testEvent struct {
	Action  string  `json:"Action"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Package string  `json:"Package"`
}

type testResult struct {
	Day  int
	Part int
	Time time.Duration
}

type yearResults struct {
	Year  int
	Tests []testResult
	Total time.Duration
}

func main() {
	var year int
	flag.IntVar(&year, "year", 0, "specific year to test (0 for all years)")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx, year, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, year int, output io.Writer) error {
	years, err := discoverAvailableYears()
	if err != nil {
		return err
	}
	if year != 0 {
		if !slices.Contains(years, year) {
			return fmt.Errorf("year %d not found; available years: %v", year, years)
		}
		years = []int{year}
	}

	results, err := collectYears(ctx, years, runTests)
	if err != nil {
		return err
	}
	return displayResults(output, results)
}

func collectYears(ctx context.Context, years []int, runYear func(context.Context, int) ([]byte, error)) ([]yearResults, error) {
	if len(years) == 0 {
		return nil, errors.New("no test years found")
	}

	results := make([]yearResults, 0, len(years))
	for _, year := range years {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		output, err := runYear(ctx, year)
		if err = errors.Join(err, ctx.Err()); err != nil {
			return nil, fmt.Errorf("year %d: %w", year, err)
		}
		tests, err := parseTestOutput(output)
		if err != nil {
			return nil, fmt.Errorf("year %d: %w", year, err)
		}
		results = append(results, aggregateResults(year, tests))
	}
	return results, nil
}

func discoverAvailableYears() ([]int, error) {
	entries, err := os.ReadDir(internalDir)
	if err != nil {
		return nil, fmt.Errorf("read internal directory: %w", err)
	}

	var years []int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		match := yearDirRegex.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}
		year, _ := strconv.Atoi(match[1])
		years = append(years, year)
	}
	slices.Sort(years)
	return years, nil
}

func testCommand(ctx context.Context, year int) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "go", "test", "-json", "-count=1", "-timeout=3m", fmt.Sprintf("./internal/aoc%d", year))
	cmd.WaitDelay = 5 * time.Second
	return cmd
}

func runTests(ctx context.Context, year int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	return runCommand(ctx, testCommand(ctx, year))
}

type boundedOutput struct {
	data     bytes.Buffer
	limit    int
	overflow bool
}

func (b *boundedOutput) Write(p []byte) (int, error) {
	n := len(p)
	remaining := b.limit - b.data.Len()
	if n > remaining {
		p = p[:remaining]
		b.overflow = true
	}
	_, _ = b.data.Write(p)
	// Continue draining the pipe so exceeding the limit cannot block the child.
	return n, nil
}

func runCommand(ctx context.Context, cmd *exec.Cmd) ([]byte, error) {
	stdout := boundedOutput{limit: maxStdout}
	stderr := boundedOutput{limit: maxStderr}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := errors.Join(cmd.Run(), ctx.Err())
	if stdout.overflow || stderr.overflow {
		err = errors.Join(err, errOutputLimit)
	}
	if err != nil {
		return stdout.data.Bytes(), fmt.Errorf("run test command: %w\nstdout:\n%s\nstderr:\n%s", err, diagnosticTail(stdout.data.Bytes()), diagnosticTail(stderr.data.Bytes()))
	}
	return stdout.data.Bytes(), nil
}

func diagnosticTail(data []byte) string {
	const limit = 8 << 10
	if len(data) > limit {
		return "[earlier output omitted]\n" + string(data[len(data)-limit:])
	}
	return string(data)
}

func parseTestOutput(output []byte) ([]testResult, error) {
	var results []testResult
	packages := make(map[string]bool)
	pending := make(map[string]bool)
	completed := make(map[string]bool)

	lineNumber := 0
	for line := range bytes.SplitSeq(output, []byte("\n")) {
		lineNumber++
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var event testEvent
		if err := json.Unmarshal(line, &event); err != nil {
			return nil, fmt.Errorf("test output line %d: %w", lineNumber, err)
		}
		if event.Action == "fail" || event.Action == "build-fail" {
			return nil, fmt.Errorf("test output line %d: %s %s %s", lineNumber, event.Package, event.Test, event.Action)
		}
		if event.Package != "" {
			if _, ok := packages[event.Package]; !ok {
				packages[event.Package] = false
			}
			if event.Test == "" && event.Action == "pass" {
				packages[event.Package] = true
			}
		}

		match := testNameRegex.FindStringSubmatch(event.Test)
		if match == nil {
			continue
		}
		switch event.Action {
		case "run":
			pending[event.Test] = true
		case "skip":
			return nil, fmt.Errorf("incomplete timing: %s skipped", event.Test)
		case "pass":
			if completed[event.Test] {
				return nil, fmt.Errorf("duplicate timing for %s", event.Test)
			}
			if event.Elapsed < 0 || event.Elapsed >= float64(math.MaxInt64)/float64(time.Second) {
				return nil, fmt.Errorf("invalid elapsed time for %s: %v", event.Test, event.Elapsed)
			}
			day, _ := strconv.Atoi(match[1])
			part, _ := strconv.Atoi(match[2])
			results = append(results, testResult{Day: day, Part: part, Time: time.Duration(event.Elapsed * float64(time.Second))})
			completed[event.Test] = true
			delete(pending, event.Test)
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no puzzle timings found")
	}
	if len(pending) != 0 {
		return nil, errors.New("incomplete timing: unfinished puzzle tests")
	}
	if len(packages) == 0 {
		return nil, errors.New("incomplete test output: missing package result")
	}
	for name, passed := range packages {
		if !passed {
			return nil, fmt.Errorf("incomplete test output: package %s did not pass", name)
		}
	}
	return results, nil
}

func aggregateResults(year int, tests []testResult) yearResults {
	var total time.Duration
	for _, test := range tests {
		total += test.Time
	}
	slices.SortFunc(tests, func(a, b testResult) int {
		if a.Day != b.Day {
			return a.Day - b.Day
		}
		return a.Part - b.Part
	})
	return yearResults{Year: year, Tests: tests, Total: total}
}

func displayResults(output io.Writer, results []yearResults) error {
	var report strings.Builder
	for _, result := range results {
		fmt.Fprintf(&report, "year %d (total: %v)\n", result.Year, result.Total.Round(time.Millisecond))
		for _, test := range result.Tests {
			fmt.Fprintf(&report, "  day%02dp%02d: %v\n", test.Day, test.Part, test.Time.Round(time.Millisecond))
		}
		report.WriteByte('\n')
	}
	if _, err := io.WriteString(output, report.String()); err != nil {
		return fmt.Errorf("write timings: %w", err)
	}
	return nil
}
