package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const (
	envSessionCookie = "AOC_SESSION_COOKIE"
	defaultTimeout   = time.Second * 30
)

func main() {
	var year, day int
	flag.IntVar(&year, "year", time.Now().Year(), "which year")
	flag.IntVar(&day, "day", time.Now().Day(), "which day")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if err := run(ctx, year, day); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, year, day int) error {
	sessionCookie := os.Getenv(envSessionCookie)
	if sessionCookie == "" {
		return fmt.Errorf("error: %s environment variable not set", envSessionCookie)
	}

	outputPath := filepath.Join("inputs", fmt.Sprintf("%d", year), fmt.Sprintf("%02d.txt", day))

	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("file already exists: %s\n", outputPath)
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error checking file: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.AddCookie(&http.Cookie{Name: "session", Value: sessionCookie})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading input: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: received status %s", resp.Status)
	}

	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
