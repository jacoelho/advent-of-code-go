package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"
)

const successfulOutput = `{"Action":"start","Package":"p"}
{"Action":"run","Package":"p","Test":"Test_day08p02"}
{"Action":"output","Package":"p","Test":"Test_day08p02","Output":"example"}
{"Action":"pass","Package":"p","Test":"Test_day08p02","Elapsed":1.25}
{"Action":"pass","Package":"p","Test":"Test_day08p02Errors","Elapsed":99}
{"Action":"pass","Package":"p","Test":"XTest_day08p02","Elapsed":88}
{"Action":"pass","Package":"p","Test":"Test_day08p02/test_00","Elapsed":77}
{"Action":"pass","Package":"p","Test":"Test_day01p01","Elapsed":0.5}
{"Action":"pass","Package":"p","Elapsed":2}
`

func TestPuzzleTimings(t *testing.T) {
	results, err := collectYears(t.Context(), []int{2024}, func(context.Context, int) ([]byte, error) {
		return []byte(successfulOutput), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Total != 1750*time.Millisecond {
		t.Fatalf("unexpected totals: %+v", results)
	}
	want := []testResult{{Day: 1, Part: 1, Time: 500 * time.Millisecond}, {Day: 8, Part: 2, Time: 1250 * time.Millisecond}}
	if !slices.Equal(results[0].Tests, want) {
		t.Fatalf("timings = %+v, want %+v", results[0].Tests, want)
	}
	var output bytes.Buffer
	if err := displayResults(&output, results); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "year 2024 (total: 1.75s)\n  day01p01: 500ms\n  day08p02: 1.25s\n\n"; got != want {
		t.Fatalf("report = %q, want %q", got, want)
	}
}

func TestIncompleteTimings(t *testing.T) {
	const passed = "{\"Action\":\"pass\",\"Package\":\"p\",\"Test\":\"Test_day01p01\",\"Elapsed\":1}\n"
	const packagePass = "{\"Action\":\"pass\",\"Package\":\"p\"}\n"
	tests := []struct {
		name   string
		output string
		cause  string
	}{
		{"malformed middle", passed + "bad json\n" + packagePass, "line 2"},
		{"truncated JSON", passed + `{"Action":`, "line 2"},
		{"missing package result", passed, "package p did not pass"},
		{"test failure", passed + "{\"Action\":\"fail\",\"Package\":\"p\",\"Test\":\"TestRegression\"}\n" + packagePass, "fail"},
		{"package failure", passed + "{\"Action\":\"fail\",\"Package\":\"p\"}\n", "fail"},
		{"skipped primary", passed + "{\"Action\":\"skip\",\"Package\":\"p\",\"Test\":\"Test_day01p02\"}\n" + packagePass, "skipped"},
		{"unfinished primary", passed + "{\"Action\":\"run\",\"Package\":\"p\",\"Test\":\"Test_day01p02\"}\n" + packagePass, "unfinished"},
		{"duplicate primary", passed + passed + packagePass, "duplicate"},
		{"no primary", packagePass, "no puzzle timings"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := parseTestOutput([]byte(tt.output))
			if err == nil || !strings.Contains(err.Error(), tt.cause) || results != nil {
				t.Fatalf("results = %v, error = %v; want no results and %q", results, err, tt.cause)
			}
		})
	}
}

func TestYearFailure(t *testing.T) {
	failure := errors.New("compiler diagnostic")
	for _, executionFailure := range []bool{false, true} {
		t.Run(fmt.Sprint(executionFailure), func(t *testing.T) {
			results, err := collectYears(t.Context(), []int{2023, 2024}, func(_ context.Context, year int) ([]byte, error) {
				if year == 2023 {
					return []byte(successfulOutput), nil
				}
				if executionFailure {
					return []byte(successfulOutput), failure
				}
				return []byte("malformed"), nil
			})
			if err == nil || !strings.Contains(err.Error(), "year 2024") || results != nil {
				t.Fatalf("results = %v, error = %v; want failed year without a report", results, err)
			}
			if executionFailure && !errors.Is(err, failure) {
				t.Fatalf("lost execution error: %v", err)
			}
		})
	}
}

func TestCanceledYear(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	results, err := collectYears(ctx, []int{2024}, func(context.Context, int) ([]byte, error) {
		cancel()
		return []byte(successfulOutput), nil
	})
	if !errors.Is(err, context.Canceled) || results != nil {
		t.Fatalf("results = %v, error = %v; want cancellation", results, err)
	}
}

func TestFreshTestCommand(t *testing.T) {
	cmd := testCommand(t.Context(), 2024)
	if len(cmd.Args) < 2 || cmd.Args[0] != "go" || cmd.Args[1] != "test" {
		t.Fatalf("unexpected command: %v", cmd.Args)
	}
	for _, arg := range []string{"-json", "-count=1", "./internal/aoc2024"} {
		if !slices.Contains(cmd.Args, arg) {
			t.Errorf("command %v lacks %s", cmd.Args, arg)
		}
	}
}

func TestCommandDiagnostics(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	cmd := helperCommand(ctx, "failure")
	output, err := runCommand(ctx, cmd)
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		t.Fatalf("error = %v, want process exit error", err)
	}
	if string(output) != "distinct stdout" || !strings.Contains(err.Error(), "distinct stdout") || !strings.Contains(err.Error(), "distinct stderr") {
		t.Fatalf("output = %q, error = %v; missing diagnostics", output, err)
	}
}

func TestCommandOutputLimit(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	output, err := runCommand(ctx, helperCommand(ctx, "overflow"))
	if !errors.Is(err, errOutputLimit) || len(output) > maxStdout {
		t.Fatalf("captured %d bytes, error = %v; want bounded output and limit error", len(output), err)
	}
}

func TestCommandCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := runCommand(ctx, helperCommand(ctx, "failure"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want cancellation", err)
	}
}

func helperCommand(ctx context.Context, mode string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestCommandHelper$")
	cmd.Env = append(os.Environ(), "AOC_TIMING_HELPER="+mode)
	return cmd
}

func TestCommandHelper(t *testing.T) {
	switch os.Getenv("AOC_TIMING_HELPER") {
	case "failure":
		fmt.Fprint(os.Stdout, "distinct stdout")
		fmt.Fprint(os.Stderr, "distinct stderr")
		os.Exit(1)
	case "overflow":
		_, _ = os.Stdout.Write(bytes.Repeat([]byte{'x'}, maxStdout+1))
		os.Exit(0)
	}
}
