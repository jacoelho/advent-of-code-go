package aoc2015

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func sumJSONNumbers(value any, ignoredValue string) (int64, error) {
	var total int64
	switch value := value.(type) {
	case json.Number:
		number, err := value.Int64()
		if err != nil {
			return 0, fmt.Errorf("invalid integer %q: %w", value, err)
		}
		return number, nil
	case []any:
		for _, item := range value {
			sum, err := sumJSONNumbers(item, ignoredValue)
			if err != nil {
				return 0, err
			}
			total += sum
		}
	case map[string]any:
		if ignoredValue != "" {
			for _, item := range value {
				if text, ok := item.(string); ok && text == ignoredValue {
					return 0, nil
				}
			}
		}
		for _, item := range value {
			sum, err := sumJSONNumbers(item, ignoredValue)
			if err != nil {
				return 0, err
			}
			total += sum
		}
	}
	return total, nil
}

func accountJSON(r io.Reader, ignoredValue string) (string, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	var document any
	if err := decoder.Decode(&document); err != nil {
		return "", err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return "", fmt.Errorf("multiple JSON documents")
		}
		return "", err
	}
	total, err := sumJSONNumbers(document, ignoredValue)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}

func day12p01(r io.Reader) (string, error) {
	return accountJSON(r, "")
}

func day12p02(r io.Reader) (string, error) {
	return accountJSON(r, "red")
}
