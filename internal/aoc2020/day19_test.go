package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const messageExample = `0: 4 1 5
1: 2 3 | 3 2
2: 4 4 | 5 5
3: 4 5 | 5 4
4: "a"
5: "b"

ababbb
bababa
abbbab
aaabbb
aaaabbb`

const recursiveMessageExample = `42: 9 14 | 10 1
9: 14 27 | 1 26
10: 23 14 | 28 1
1: "a"
11: 42 31
5: 1 14 | 15 1
19: 14 1 | 14 14
12: 24 14 | 19 1
16: 15 1 | 14 14
31: 14 17 | 1 13
6: 14 14 | 1 14
2: 1 24 | 14 4
0: 8 11
13: 14 3 | 1 12
15: 1 | 14
17: 14 2 | 1 7
23: 25 1 | 22 14
28: 16 1
4: 1 1
20: 14 14 | 1 15
3: 5 14 | 16 1
27: 1 6 | 14 18
14: "b"
21: 14 1 | 1 14
25: 1 1 | 1 14
22: 14 14
8: 42
26: 14 22 | 1 20
18: 15 15
7: 14 5 | 1 21
24: 14 1

abbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa
bbabbbbaabaabba
babbbbaabbbbbabbbbbbaabaaabaaa
aaabbbbbbaaaabaababaabababbabaaabbababababaaa
bbbbbbbaaaabbbbaaabbabaaa
bbbababbbbaaaaaaaabbababaaababaabab
ababaaaaaabaaab
ababaaaaabbbaba
baabbaaaabbaaaababbaababb
abbbbabbbbaaaababbbbbbaaaababb
aaaaabbaabaaaaababaa
aaaabbaaaabbaaa
aaaabbaabbaaaaaaabbbabbbaaabbaabaaa
babaaabbbaaabaababbaabababaaab
aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba`

func Test_day19p01(t *testing.T) {
	aoc.AOCTest(t, day19p01, []aoc.TestInput{
		{Input: strings.NewReader(messageExample), Want: "2"},
		{Input: strings.NewReader(recursiveMessageExample), Want: "3"},
		{Input: aoc.FileInput(t, 2020, 19), Want: "162"},
	})
}

func Test_day19p02(t *testing.T) {
	aoc.AOCTest(t, day19p02, []aoc.TestInput{
		{Input: strings.NewReader(recursiveMessageExample), Want: "12"},
		{Input: aoc.FileInput(t, 2020, 19), Want: "267"},
	})
}

func TestMessageGrammarConsumesWholeMessage(t *testing.T) {
	aoc.AOCTest(t, day19p01, []aoc.TestInput{
		{Input: strings.NewReader("0: 1 | 1 1\n1: \"a\"\n\na\naa\naaa\nb"), Want: "2"},
		{Input: strings.NewReader("0: 1 | 1 1 | 1 1 1\n1: \"a\"\n\na\naa\naaa\naaaa"), Want: "3"},
		{Input: strings.NewReader("0: 1 | 2\n1: \"a\"\n2: \"a\"\n\na"), Want: "1"},
	})
}

func TestRecursiveMessageGrammar(t *testing.T) {
	const grammar = "0: 8 11\n8: 42\n11: 42 31\n42: \"a\"\n31: \"b\"\n\n"
	messages := "aab\naaab\nab\naabb\n" + strings.Repeat("a", 32) + strings.Repeat("b", 31)
	aoc.AOCTest(t, day19p02, []aoc.TestInput{{Input: strings.NewReader(grammar + messages), Want: "3"}})
	aoc.AOCTest(t, day19p01, []aoc.TestInput{
		{Input: strings.NewReader("0: 0 1 | 1\n1: \"a\"\n\na\naa\naaa\nb"), Want: "3"},
		{Input: strings.NewReader("0: 1\n1: 0 | 2\n2: \"a\"\n\na\naa"), Want: "1"},
		{Input: strings.NewReader("0: 0\n\na"), Want: "0"},
	})
}

func TestMessageGrammarErrors(t *testing.T) {
	for _, input := range []string{"0: 1\n\na", "1: \"a\"\n\na", "0: 1 |\n1: \"a\"\n\na", "0: \"ab\"\n\na", "0: \"a\"\n0: \"b\"\n\na"} {
		if _, err := day19p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected grammar error for %q", input)
		}
	}
	for _, solve := range []func(io.Reader) (string, error){day19p01, day19p02} {
		failure := errors.New("read messages")
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("got %v, want %v", err, failure)
		}
	}
}
