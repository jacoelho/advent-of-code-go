package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

const (
	internalDir = "internal"
	passAction  = "pass"
)

var (
	yearDirRegex  = regexp.MustCompile(`^aoc(\d{4})$`)
	testNameRegex = regexp.MustCompile(`Test_day(\d{2})p(\d{2})`)
)

func testKey(day, part int) string {
	return fmt.Sprintf("day%02dp%02d", day, part)
}

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

	if err := run(ctx, year); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, year int) error {
	availableYears, err := discoverAvailableYears()
	if err != nil {
		return fmt.Errorf("error discovering available years: %w", err)
	}

	var yearsToTest []int
	if year != 0 {
		if !slices.Contains(availableYears, year) {
			return fmt.Errorf("year %d not found. Available years: %v", year, availableYears)
		}
		yearsToTest = []int{year}
	} else {
		yearsToTest = availableYears
	}

	var allResults []yearResults

	for _, y := range yearsToTest {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		output, err := runTests(ctx, y)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error testing year %d: %v\n", y, err)
			continue
		}

		tests, err := parseTestOutput(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing tests for year %d: %v\n", y, err)
			continue
		}

		if len(tests) > 0 {
			results := aggregateResults(y, tests)
			allResults = append(allResults, results)
		}
	}

	displayResults(allResults)
	return nil
}

func discoverAvailableYears() ([]int, error) {
	entries, err := os.ReadDir(internalDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read internal directory: %w", err)
	}

	dirs := xslices.Filter(func(entry os.DirEntry) bool {
		return entry.IsDir()
	}, entries)

	years := xslices.Map(func(entry os.DirEntry) int {
		matches := yearDirRegex.FindStringSubmatch(entry.Name())
		if len(matches) == 2 {
			year, _ := strconv.Atoi(matches[1])
			return year
		}
		return 0
	}, dirs)

	validYears := xslices.Filter(func(year int) bool {
		return year > 0
	}, years)

	sort.Ints(validYears)
	return validYears, nil
}

func parseTestEvent(event testEvent) (testResult, bool) {
	if event.Action != passAction {
		return testResult{}, false
	}

	matches := testNameRegex.FindStringSubmatch(event.Test)
	if len(matches) != 3 {
		return testResult{}, false
	}

	day, _ := strconv.Atoi(matches[1])
	part, _ := strconv.Atoi(matches[2])
	duration := time.Duration(event.Elapsed * float64(time.Second))

	return testResult{
		Day:  day,
		Part: part,
		Time: duration,
	}, true
}

func runTests(ctx context.Context, year int) ([]byte, error) {
	packagePath := fmt.Sprintf("./internal/aoc%d/...", year)

	cmd := exec.CommandContext(ctx, "go", "test", "-json", packagePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run tests: %w", err)
	}

	return output, nil
}

func parseTestOutput(output []byte) ([]testResult, error) {
	testResults := make(map[string]testResult)

	lines := strings.SplitSeq(string(output), "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		var event testEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		if result, ok := parseTestEvent(event); ok {
			testResults[testKey(result.Day, result.Part)] = result
		}
	}

	return slices.Collect(maps.Values(testResults)), nil
}

func aggregateResults(year int, tests []testResult) yearResults {
	totalTime := xslices.Sum(xslices.Map(func(test testResult) time.Duration {
		return test.Time
	}, tests))

	slices.SortFunc(tests, func(a, b testResult) int {
		if a.Day != b.Day {
			return a.Day - b.Day
		}
		return a.Part - b.Part
	})

	return yearResults{
		Year:  year,
		Tests: tests,
		Total: totalTime,
	}
}

func displayResults(results []yearResults) {
	for _, yearResult := range results {
		fmt.Printf("year %d (total: %v)\n", yearResult.Year, yearResult.Total.Round(time.Millisecond))

		for _, test := range yearResult.Tests {
			fmt.Printf("  day%02dp%02d: %v\n", test.Day, test.Part, test.Time.Round(time.Millisecond))
		}
		fmt.Println()
	}
}
