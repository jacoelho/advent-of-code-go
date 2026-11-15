package aoc2023

import (
	"fmt"
	"strconv"
	"strings"
)

// parseNaturals parses the whitespace-separated nonnegative integers in s.
func parseNaturals(s string) ([]int, error) {
	fields := strings.Fields(s)
	values := make([]int, len(fields))
	for i, field := range fields {
		value, err := strconv.Atoi(field)
		if err != nil || value < 0 {
			return nil, fmt.Errorf("invalid nonnegative integer %q", field)
		}
		values[i] = value
	}
	return values, nil
}
