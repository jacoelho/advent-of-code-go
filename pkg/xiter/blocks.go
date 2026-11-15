package xiter

import "iter"

// Blocks groups lines into blocks separated by blank lines. Consecutive blank
// lines enclose an empty block and a leading blank line yields an empty first
// block, but a single trailing blank line only ends the last block. An empty
// source yields no blocks. Each block is a new slice that the caller owns.
func Blocks(lines iter.Seq[string]) iter.Seq[[]string] {
	return func(yield func([]string) bool) {
		var block []string
		for line := range lines {
			if line != "" {
				block = append(block, line)
				continue
			}
			if !yield(block) {
				return
			}
			block = nil
		}
		if len(block) > 0 {
			yield(block)
		}
	}
}
