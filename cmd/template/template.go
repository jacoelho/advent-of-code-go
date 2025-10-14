package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type config struct {
	Year string
	Day  string
}

var dayTmpl = template.Must(template.New("day").Parse(`package aoc{{ .Year }}

import "io"

func day{{ .Day }}p01(r io.Reader) (string, error) {
	return "", nil
}

func day{{ .Day }}p02(r io.Reader) (string, error) {
	return "", nil
}`))

var testTmpl = template.Must(template.New("test").Parse(`package aoc{{ .Year }}

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day{{ .Day }}p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(` + "``" + `),
			Want: "",
		},
		{
			Input: aoc.FileInput(t, {{ .Year }}, {{ .Day }}),
			Want:  "",
		},
	}
	aoc.AOCTest(t, day{{ .Day }}p01, tests)
}

func Test_day{{ .Day }}p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(` + "``" + `),
			Want: "",
		},
		{
			Input: aoc.FileInput(t, {{ .Year }}, {{ .Day }}),
			Want:  "",
		},
	}
	aoc.AOCTest(t, day{{ .Day }}p02, tests)
}`))

func main() {
	var c config
	flag.StringVar(&c.Year, "year", time.Now().Format("2006"), "which year")
	flag.StringVar(&c.Day, "day", time.Now().Format("02"), "which day")
	flag.Parse()

	baseFile := filepath.Join("internal", "aoc"+c.Year, "day"+c.Day)
	dayFile := baseFile + ".go"
	testFile := baseFile + "_test.go"

	for _, f := range []string{dayFile, testFile} {
		if _, err := os.Stat(f); err == nil || !errors.Is(err, os.ErrNotExist) {
			panic("file already exists " + f)
		}
	}

	f, err := os.OpenFile(baseFile+".go", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := dayTmpl.Execute(f, c); err != nil {
		panic(err)
	}

	tst, err := os.OpenFile(baseFile+"_test.go", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer tst.Close()

	if err := testTmpl.Execute(tst, c); err != nil {
		panic(err)
	}
}
