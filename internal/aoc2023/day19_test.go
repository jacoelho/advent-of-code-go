package aoc2023

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day19p01(t *testing.T) {
	exampleInput := `px{a<2006:qkq,m>2090:A,rfg}
pv{a>1716:R,A}
lnx{m>1548:A,A}
rfg{s<537:gd,x>2440:R,A}
qs{s>3448:A,lnx}
qkq{x<1416:A,crn}
crn{x>2662:A,R}
in{s<1351:px,qqz}
qqz{s>2770:qs,m<1801:hdj,R}
gd{a>3333:R,R}
hdj{m>838:A,pv}

{x=787,m=2655,a=1222,s=2876}
{x=1679,m=44,a=2067,s=496}
{x=2036,m=264,a=79,s=2244}
{x=2461,m=1339,a=466,s=291}
{x=2127,m=1623,a=2188,s=1013}`

	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(exampleInput),
			Want:  "19114",
		},
		{
			Input: aoc.FileInput(t, 2023, 19),
			Want:  "368964",
		},
	}
	aoc.AOCTest(t, day19p01, tests)
}

func Test_day19p02(t *testing.T) {
	exampleInput := `px{a<2006:qkq,m>2090:A,rfg}
pv{a>1716:R,A}
lnx{m>1548:A,A}
rfg{s<537:gd,x>2440:R,A}
qs{s>3448:A,lnx}
qkq{x<1416:A,crn}
crn{x>2662:A,R}
in{s<1351:px,qqz}
qqz{s>2770:qs,m<1801:hdj,R}
gd{a>3333:R,R}
hdj{m>838:A,pv}

{x=787,m=2655,a=1222,s=2876}
{x=1679,m=44,a=2067,s=496}
{x=2036,m=264,a=79,s=2244}
{x=2461,m=1339,a=466,s=291}
{x=2127,m=1623,a=2188,s=1013}`

	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(exampleInput),
			Want:  "167409079868000",
		},
		{
			Input: aoc.FileInput(t, 2023, 19),
			Want:  "127675188176682",
		},
	}
	aoc.AOCTest(t, day19p02, tests)
}

func TestWorkflowRangeBounds(t *testing.T) {
	const part = "\n\n{x=1,m=1,a=1,s=1}"
	tests := []struct{ workflows, want string }{
		{"in{x<10:small,R}\nsmall{x<20:A,R}", "576000000000"},
		{"in{x>3991:large,R}\nlarge{x>3981:A,R}", "576000000000"},
		{"in{x<10:small,R}\nsmall{x>20:R,A}", "576000000000"},
		{"in{x<1:A,R}", "0"},
		{"in{x>4000:A,R}", "0"},
		{"in{x<4001:A,R}", "256000000000000"},
		{"in{x<-9223372036854775808:R,A}", "256000000000000"},
		{"in{x>9223372036854775807:A,R}", "0"},
		{"in{x<9223372036854775807:A,R}", "256000000000000"},
		{"in{x>-9223372036854775808:A,R}", "256000000000000"},
	}
	for _, tt := range tests {
		aoc.AOCTest(t, day19p02, []aoc.TestInput{{Input: strings.NewReader(tt.workflows + part), Want: tt.want}})
	}
}

func TestWorkflowTinyDomain(t *testing.T) {
	const definitions = "in{x<3:left,m>3:A,R}\nleft{a>2:A,s<3:A,R}"
	workflows, err := parseWorkflows(definitions)
	if err != nil {
		t.Fatal(err)
	}
	var want int64
	for x := 1; x <= 4; x++ {
		for m := 1; m <= 4; m++ {
			for a := 1; a <= 4; a++ {
				for s := 1; s <= 4; s++ {
					accepted := m > 3
					if x < 3 {
						accepted = a > 2 || s < 3
					}
					if accepted {
						want++
					}
				}
			}
		}
	}
	got, err := countAccepted(workflows, 4)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %d accepted ratings, want %d", got, want)
	}
	const parts = "\n\n{x=1,m=1,a=3,s=4}\n{x=2,m=1,a=2,s=2}\n{x=3,m=4,a=1,s=1}\n{x=4,m=3,a=4,s=4}"
	aoc.AOCTest(t, day19p01, []aoc.TestInput{{Input: strings.NewReader(definitions + parts), Want: "25"}})
}

func TestWorkflowCycles(t *testing.T) {
	const dead = "in{x<1:loop,A}\nloop{loop}\n\n{x=1,m=1,a=1,s=1}"
	aoc.AOCTest(t, day19p01, []aoc.TestInput{{Input: strings.NewReader(dead), Want: "4"}})
	aoc.AOCTest(t, day19p02, []aoc.TestInput{{Input: strings.NewReader(dead), Want: "256000000000000"}})
	const reachable = "in{x<2:loop,A}\nloop{loop}"
	aoc.AOCTest(t, day19p01, []aoc.TestInput{{Input: strings.NewReader(reachable + "\n\n{x=2,m=1,a=1,s=1}"), Want: "5"}})
	for _, solve := range []func(io.Reader) (string, error){day19p01, day19p02} {
		if _, err := solve(strings.NewReader(reachable + "\n\n{x=1,m=1,a=1,s=1}")); err == nil {
			t.Fatal("expected reachable workflow cycle error")
		}
	}
}

func TestWorkflowInputErrors(t *testing.T) {
	for _, input := range []string{
		"in{A}\n\n{x=1,x=2,a=1,s=1}",
		"in{A}\n\n{x=1,m=1,a=1}",
		"in{A}\n\n{x=0,m=1,a=1,s=1}",
		"in{missing}\n\n{x=1,m=1,a=1,s=1}",
		"in{x<2:A}\n\n{x=1,m=1,a=1,s=1}",
		"in{A}\nin{R}\n\n{x=1,m=1,a=1,s=1}",
	} {
		for _, solve := range []func(io.Reader) (string, error){day19p01, day19p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected workflow error for %q", input)
			}
		}
	}
	failure := errors.New("workflow read failed")
	for _, solve := range []func(io.Reader) (string, error){day19p01, day19p02} {
		_, err := solve(io.MultiReader(strings.NewReader("in{A}\n\n{x=1,m=1,a=1,s=1}"), iotest.ErrReader(failure)))
		if !errors.Is(err, failure) {
			t.Fatalf("got %v, want reader failure", err)
		}
	}
}
