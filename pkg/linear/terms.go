package linear

import "cmp"

type term[K Key] struct {
	variable    Var[K]
	coefficient Rat
}

// mergeTerms requires key-sorted inputs. The result may share an input list;
// see the Expr.terms invariant.
func mergeTerms[K Key](left, right []term[K]) []term[K] {
	if len(left) == 0 {
		return right
	}
	if len(right) == 0 {
		return left
	}
	result := make([]term[K], 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		a, b := left[i], right[j]
		switch cmp.Compare(a.variable.key, b.variable.key) {
		case -1:
			result = append(result, a)
			i++
		case 1:
			result = append(result, b)
			j++
		case 0:
			if coefficient := a.coefficient.Add(b.coefficient); !coefficient.IsZero() {
				result = append(result, term[K]{variable: a.variable, coefficient: coefficient})
			}
			i++
			j++
		}
	}
	result = append(result, left[i:]...)
	return append(result, right[j:]...)
}
