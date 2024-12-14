package convert

import (
	"fmt"
	"reflect"
	"testing"
)

func TestToDigits(t *testing.T) {
	type testCase struct {
		input int
		want  []int
	}
	tests := []testCase{
		{
			input: 0,
			want:  []int{0},
		},
		{
			input: 1,
			want:  []int{1},
		},
		{
			input: 54321,
			want:  []int{5, 4, 3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			if got := ToDigits(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractDigits(t *testing.T) {
	type testCase struct {
		input string
		want  []int
	}
	tests := []testCase{
		{
			input: "123 -456 78 -9abc-10,0 - 3",
			want:  []int{123, -456, 78, -9, -10, 0, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ExtractDigits[int](tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}
