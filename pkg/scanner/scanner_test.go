package scanner

import (
	"errors"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"
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

func TestConversionFailureIsTerminal(t *testing.T) {
	s := NewScanner(strings.NewReader("1\nbad\n3"), func(v []byte) (int, error) {
		return strconv.Atoi(string(v))
	})
	values := s.Values()
	if got := slices.Collect(values); !slices.Equal(got, []int{1}) {
		t.Fatalf("prefix = %v, want [1]", got)
	}
	err := s.Err()
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("Err() = %v, want conversion syntax error", err)
	}
	if got := slices.Collect(values); len(got) != 0 || s.Err() != err {
		t.Fatalf("resumed failed scanner returned %v, error %v; want no values and %v", got, s.Err(), err)
	}
}

func TestSuccessfulEarlyStopCanResume(t *testing.T) {
	s := NewScanner(strings.NewReader("1\n2\n3"), func(v []byte) (int, error) {
		return strconv.Atoi(string(v))
	})
	for value := range s.Values() {
		if value != 1 {
			t.Fatalf("first value = %d, want 1", value)
		}
		break
	}
	if got := slices.Collect(s.Values()); !slices.Equal(got, []int{2, 3}) || s.Err() != nil {
		t.Fatalf("resumed scanner = %v, error %v; want [2 3], nil", got, s.Err())
	}
}

func TestReaderErrorAfterValidPrefix(t *testing.T) {
	wantErr := errors.New("reader failed")
	r := io.MultiReader(strings.NewReader("1\n2\n"), iotest.ErrReader(wantErr))
	s := NewScanner(r, func(v []byte) (int, error) { return strconv.Atoi(string(v)) })
	if got := slices.Collect(s.Values()); !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("prefix = %v, want [1 2]", got)
	}
	if !errors.Is(s.Err(), wantErr) {
		t.Fatalf("Err() = %v, want %v", s.Err(), wantErr)
	}
}
