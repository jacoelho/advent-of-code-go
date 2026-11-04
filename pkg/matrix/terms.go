package matrix

import "cmp"

type integerTerm[K Key] struct {
	variable    Var[K]
	coefficient int64
}

// termSet keeps sorted, unique, nonzero terms. Published storage is immutable.
// Inline storage holds at most two compact integer coefficients.
// Zero inline coefficients mark unused slots.
// Spill storage contains the entire sequence; it never points into inline.
type termSet[K Key] struct {
	inline [2]integerTerm[K]
	spill  []term[K]
}

func singleIntegerTerm[K Key](variable Var[K], coefficient int64) termSet[K] {
	if coefficient == 0 {
		return termSet[K]{}
	}
	return termSet[K]{
		inline: [2]integerTerm[K]{{variable: variable, coefficient: coefficient}},
	}
}

func singleTerm[K Key](variable Var[K], coefficient Rat) termSet[K] {
	if coefficient.value == nil && coefficient.denominator == 0 {
		return singleIntegerTerm(variable, coefficient.numerator)
	}
	return termSet[K]{
		spill: []term[K]{{variable: variable, coefficient: coefficient}},
	}
}

func (s *termSet[K]) len() int {
	if s.spill != nil {
		return len(s.spill)
	}
	if s.inline[1].coefficient != 0 {
		return 2
	}
	if s.inline[0].coefficient != 0 {
		return 1
	}
	return 0
}

func (s *termSet[K]) at(index int) term[K] {
	if s.spill != nil {
		return s.spill[index]
	}
	value := s.inline[index]
	return term[K]{variable: value.variable, coefficient: integerRat(value.coefficient)}
}

// append is only valid on a fresh, unpublished builder, in canonical term order.
func (s *termSet[K]) append(value term[K]) {
	count := s.len()
	if s.spill != nil {
		s.spill = append(s.spill, value)
	} else if count < len(s.inline) && value.coefficient.value == nil && value.coefficient.denominator == 0 {
		s.inline[count] = integerTerm[K]{variable: value.variable, coefficient: value.coefficient.numerator}
	} else {
		spill := make([]term[K], count+1)
		for i := range count {
			spill[i] = s.at(i)
		}
		spill[count] = value
		s.spill = spill
		s.inline = [2]integerTerm[K]{}
	}
}

func mergeTerms[K Key](left, right *termSet[K]) termSet[K] {
	var result termSet[K]
	leftCount, rightCount := left.len(), right.len()
	i, j := 0, 0
	for i < leftCount && j < rightCount {
		a, b := left.at(i), right.at(j)
		switch cmp.Compare(a.variable.key, b.variable.key) {
		case -1:
			result.append(a)
			i++
		case 1:
			result.append(b)
			j++
		case 0:
			coefficient := a.coefficient.Add(b.coefficient)
			if !coefficient.IsZero() {
				result.append(term[K]{variable: a.variable, coefficient: coefficient})
			}
			i++
			j++
		}
	}
	for ; i < leftCount; i++ {
		result.append(left.at(i))
	}
	for ; j < rightCount; j++ {
		result.append(right.at(j))
	}
	return result
}
