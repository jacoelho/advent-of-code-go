package scanner

import (
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	convert := func(v []byte) (int, error) {
		return strconv.Atoi(string(v))
	}

	type testCase struct {
		name    string
		input   io.Reader
		want    []int
		wantErr bool
	}
	tests := []testCase{
		{
			name:  "empty",
			input: strings.NewReader(""),
			want:  nil,
		},
		{
			name:  "regular case",
			input: strings.NewReader("1\n2\n3"),
			want:  []int{1, 2, 3},
		},
		{
			name:  "happy case \n terminated",
			input: strings.NewReader("1\n2\n3\n"),
			want:  []int{1, 2, 3},
		},
		{
			name:    "convert error",
			input:   strings.NewReader("1\na\n3"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := NewScanner(tt.input, convert)

			got := slices.Collect(lines.Values())
			if tt.wantErr != (lines.Err() != nil) {
				t.Fatalf("unexpected error: %v", lines.Err())
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}
