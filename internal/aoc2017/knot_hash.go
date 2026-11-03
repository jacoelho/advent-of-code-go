package aoc2017

func knotRound(values []byte, lengths []int, position, skip *int) {
	for _, length := range lengths {
		for i := 0; i < length/2; i++ {
			left, right := (*position+i)%len(values), (*position+length-1-i)%len(values)
			values[left], values[right] = values[right], values[left]
		}
		*position = (*position + length + *skip) % len(values)
		*skip++
	}
}

func knotHash(input []byte) [16]byte {
	lengths := make([]int, 0, len(input)+5)
	for _, value := range input {
		lengths = append(lengths, int(value))
	}
	lengths = append(lengths, 17, 31, 73, 47, 23)
	var values [256]byte
	for i := range values {
		values[i] = byte(i)
	}

	position, skip := 0, 0
	for range 64 {
		knotRound(values[:], lengths, &position, &skip)
	}

	var hash [16]byte
	for i, value := range values {
		hash[i/16] ^= value
	}
	return hash
}
