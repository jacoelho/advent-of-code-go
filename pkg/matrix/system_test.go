package matrix_test

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/pkg/matrix"
)

func ExampleSystem() {
	a, b := matrix.Variable("a"), matrix.Variable("b")
	solution, err := matrix.NewSystem(a, b).
		Add(a.Scale(2).Add(b.Scale(1)).Equal(4)).
		Solve()
	if err != nil {
		fmt.Println(err)
		return
	}
	values, err := solution.Bind(map[matrix.Var[string]]matrix.Rat{
		b: matrix.NewRat(2, 1),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("a =", values[a].Numerator())
	fmt.Println("b =", values[b].Numerator())
	// Output:
	// a = 1
	// b = 2
}

func solve[K matrix.Key](t *testing.T, system *matrix.System[K]) *matrix.Solution[K] {
	t.Helper()
	solution, err := system.Solve()
	if err != nil {
		t.Fatal(err)
	}
	return solution
}

func wantRat(t *testing.T, got matrix.Rat, numerator, denominator int64) {
	t.Helper()
	if got.Cmp(matrix.NewRat(numerator, denominator)) != 0 {
		t.Fatalf("got %s/%s; want %d/%d", got.Numerator(), got.Denominator(), numerator, denominator)
	}
}

func TestSystemUnique(t *testing.T) {
	type component uint8
	type coefficient int64
	x, y := matrix.Variable(component(0)), matrix.Variable(component(1))
	solution := solve(t, matrix.NewSystem(x, y).
		Add(y.Scale(coefficient(2)).Equal(int64(8))).
		Add(x.Scale(3).Add(y.Scale(1)).Equal(7)))
	if len(solution.Free()) != 0 {
		t.Fatal("unique system has free variables")
	}
	values, err := solution.Bind(nil)
	if err != nil || len(values) != 2 {
		t.Fatalf("Bind(nil) = %v, %v", values, err)
	}
	wantRat(t, values[x], 1, 1)
	wantRat(t, values[y], 4, 1)
	if x.Key() != 0 || y.Key() != 1 {
		t.Fatal("variable keys changed")
	}
	captured, err := solution.Expression(x)
	if err != nil {
		t.Fatal(err)
	}
	*solution = matrix.Solution[component]{}
	wantRat(t, captured.Constant(), 1, 1)
}

func TestSystemRationalFamily(t *testing.T) {
	a, b := matrix.Variable("a"), matrix.Variable("b")
	solution := solve(t, matrix.NewSystem(a, b).
		Add(a.Scale(2).Add(b.Scale(1)).Equal(4)).
		Add(a.Scale(4).Add(b.Scale(2)).Equal(8)))
	if !slices.Equal(solution.Free(), []matrix.Var[string]{b}) {
		t.Fatal("expected b to be free")
	}
	expression, err := solution.Expression(matrix.Variable("a"))
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, expression.Constant(), 2, 1)
	wantRat(t, expression.Coefficient(b), -1, 2)
	wantRat(t, expression.Coefficient(a), 0, 1)
	for _, binding := range []struct{ b, numerator, denominator int64 }{{1, 3, 2}, {2, 1, 1}} {
		values, err := solution.Bind(map[matrix.Var[string]]matrix.Rat{b: matrix.NewRat(binding.b, 1)})
		if err != nil {
			t.Fatal(err)
		}
		wantRat(t, values[a], binding.numerator, binding.denominator)
		wantRat(t, values[b], binding.b, 1)
	}
	// A returned affine expression can become another equation, including its constant.
	bound := solve(t, matrix.NewSystem(b).Add(expression.Equal(1)))
	values, err := bound.Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[b], 2, 1)

	exact := solve(t, matrix.NewSystem(a).
		Add(a.ScaleRat(matrix.NewRat(1, 2)).EqualRat(matrix.NewRat(3, 2))))
	values, err = exact.Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 3, 1)
}

func TestAffineExpressionOperationsAndSnapshots(t *testing.T) {
	x, z, parameter := matrix.Variable("x"), matrix.Variable("z"), matrix.Variable("t")
	solution := solve(t, matrix.NewSystem(x, z, parameter).
		Add(x.Scale(1).Add(parameter.Scale(1)).Equal(2)).
		Add(x.Scale(1).Add(z.Scale(1)).Equal(0)))
	expressionX, err := solution.Expression(x)
	if err != nil {
		t.Fatal(err)
	}
	expressionZ, err := solution.Expression(z)
	if err != nil {
		t.Fatal(err)
	}
	captured := expressionX.Equal(1)
	doubled := expressionX.Add(expressionX)
	canceled := expressionX.Add(expressionZ)
	for _, tt := range []struct {
		name                  string
		expression            matrix.Expr[string]
		constant, coefficient int64
	}{
		{"original x", expressionX, 2, -1},
		{"original z", expressionZ, -2, 1},
		{"doubled", doubled, 4, -2},
		{"canceled", canceled, 0, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			wantRat(t, tt.expression.Constant(), tt.constant, 1)
			wantRat(t, tt.expression.Coefficient(parameter), tt.coefficient, 1)
			wantRat(t, tt.expression.Coefficient(x), 0, 1)
			wantRat(t, tt.expression.Coefficient(z), 0, 1)
		})
	}
	values, err := solve(t, matrix.NewSystem(parameter).Add(captured)).Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[parameter], 1, 1)
	values, err = solve(t, matrix.NewSystem[string]().Add(canceled.Equal(0))).Bind(nil)
	if err != nil || len(values) != 0 {
		t.Fatalf("canceled affine expression = %v, %v", values, err)
	}
}

func TestSystemTwoVariableExactSolutions(t *testing.T) {
	x, y := matrix.Variable("x"), matrix.Variable("y")
	for _, tt := range []struct {
		name   string
		first  matrix.Equation[string]
		second matrix.Equation[string]
		wantX  matrix.Rat
		wantY  matrix.Rat
	}{
		{
			"fractional signed solution",
			x.Scale(2).Add(y.Scale(1)).Equal(0),
			x.Scale(1).Add(y.Scale(2)).Equal(1),
			matrix.NewRat(-1, 3), matrix.NewRat(2, 3),
		},
		{
			"fractional coefficient",
			x.ScaleRat(matrix.NewRat(1, 2)).Add(y.Scale(1)).Equal(2),
			x.Scale(1).Add(y.Scale(3)).Equal(5),
			matrix.NewRat(2, 1), matrix.NewRat(1, 1),
		},
		{
			"fractional right hand side",
			x.Scale(2).Add(y.Scale(1)).EqualRat(matrix.NewRat(1, 2)),
			x.Scale(1).Add(y.Scale(2)).Equal(1),
			matrix.NewRat(0, 1), matrix.NewRat(1, 2),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			solution := solve(t, matrix.NewSystem(x, y).Add(tt.first).Add(tt.second))
			values, err := solution.Bind(nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(values) != 2 || values[x].Cmp(tt.wantX) != 0 || values[y].Cmp(tt.wantY) != 0 {
				t.Fatalf("unexpected exact solution: %v", values)
			}
		})
	}
}

func TestSystemArbitraryPrecisionCoefficients(t *testing.T) {
	x, y := matrix.Variable("x"), matrix.Variable("y")
	large := new(big.Int).Lsh(big.NewInt(1), 80)
	coefficient := matrix.NewInteger(large)
	system := matrix.NewSystem(x, y).
		Add(x.ScaleRat(coefficient).EqualRat(matrix.NewInteger(new(big.Int).Mul(large, big.NewInt(2))))).
		Add(y.ScaleRat(coefficient).EqualRat(matrix.NewInteger(new(big.Int).Mul(large, big.NewInt(3)))))
	for range 2 {
		values, err := solve(t, system).Bind(nil)
		if err != nil {
			t.Fatal(err)
		}
		wantRat(t, values[x], 2, 1)
		wantRat(t, values[y], 3, 1)
	}
	if coefficient.Numerator().Cmp(large) != 0 {
		t.Fatal("solving mutated an input coefficient")
	}
}

func TestSystemTwoVariableIntegerOverflow(t *testing.T) {
	x, y := matrix.Variable("x"), matrix.Variable("y")
	for _, tt := range []struct {
		name          string
		first, second matrix.Equation[string]
		wantX, wantY  string
	}{
		{"determinant subtraction", x.Scale(int64(math.MaxInt64)).Add(y.Scale(1)).Equal(1),
			x.Scale(int64(math.MinInt64)).Add(y.Scale(1)).Equal(1), "0", "1"},
		{"first numerator product", x.Scale(1).Add(y.Scale(2)).Equal(0),
			y.Scale(1).Equal(int64(math.MaxInt64)), "-18446744073709551614", "9223372036854775807"},
		{"second numerator product", x.Scale(1).Equal(int64(math.MaxInt64)),
			x.Scale(2).Add(y.Scale(1)).Equal(0), "9223372036854775807", "-18446744073709551614"},
		{"numerator subtraction", x.Scale(1).Add(y.Scale(-1)).Equal(int64(math.MaxInt64)),
			y.Scale(1).Equal(1), "9223372036854775808", "1"},
		{"minimum divided by negative one", x.Scale(-1).Equal(int64(math.MinInt64)),
			y.Scale(1).Equal(0), "9223372036854775808", "0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			values, err := solve(t, matrix.NewSystem(x, y).Add(tt.first).Add(tt.second)).Bind(nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []struct {
				variable matrix.Var[string]
				integer  string
			}{{x, tt.wantX}, {y, tt.wantY}} {
				got := values[want.variable]
				if got.Numerator().String() != want.integer || got.Denominator().String() != "1" {
					t.Fatalf("%s = %s/%s; want %s", want.variable.Key(), got.Numerator(), got.Denominator(), want.integer)
				}
			}
		})
	}
}

func TestSystemMultipleFreeVariables(t *testing.T) {
	a, b, c, d := matrix.Variable("a"), matrix.Variable("b"), matrix.Variable("c"), matrix.Variable("d")
	system := matrix.NewSystem(c, a, d, b).
		Add(d.Scale(3).Add(a.Scale(1)).Add(b.Scale(2)).Equal(7)).
		Add(c.Scale(1).Add(d.Scale(-4)).Equal(5))
	for _, solution := range []*matrix.Solution[string]{solve(t, system), solve(t, system)} {
		if !slices.Equal(solution.Free(), []matrix.Var[string]{d, b}) {
			t.Fatal("expected free variables d, b in declaration order")
		}
		expression, err := solution.Expression(a)
		if err != nil {
			t.Fatal(err)
		}
		wantRat(t, expression.Constant(), 7, 1)
		wantRat(t, expression.Coefficient(b), -2, 1)
		wantRat(t, expression.Coefficient(d), -3, 1)
		values, err := solution.Bind(map[matrix.Var[string]]matrix.Rat{b: matrix.NewRat(1, 1), d: matrix.NewRat(2, 1)})
		if err != nil || len(values) != 4 {
			t.Fatalf("Bind = %v, %v", values, err)
		}
		wantRat(t, values[a], -1, 1)
		wantRat(t, values[b], 1, 1)
		wantRat(t, values[c], 13, 1)
		wantRat(t, values[d], 2, 1)
	}
}

func TestExpressionMergingAndCancellation(t *testing.T) {
	a, b, c, d := matrix.Variable("a"), matrix.Variable("b"), matrix.Variable("c"), matrix.Variable("d")
	left := d.Scale(4).Add(b.Scale(2)).Add(a.Scale(1))
	right := c.Scale(3).Add(b.Scale(-2)).Add(a.Scale(5))
	merged := left.Add(right)
	for _, tt := range []struct {
		expression matrix.Expr[string]
		want       [4]int64
	}{
		{left, [4]int64{1, 2, 0, 4}},
		{right, [4]int64{5, -2, 3, 0}},
		{merged, [4]int64{6, 0, 3, 4}},
	} {
		for i, variable := range []matrix.Var[string]{a, b, c, d} {
			wantRat(t, tt.expression.Coefficient(variable), tt.want[i], 1)
		}
	}
	captured := merged.Equal(23)
	canceled := merged.Add(a.Scale(-6)).Add(c.Scale(-3)).Add(d.Scale(-4))
	values, err := solve(t, matrix.NewSystem[string]().Add(canceled.Equal(0))).Bind(nil)
	if err != nil || len(values) != 0 {
		t.Fatalf("canceled expression = %v, %v", values, err)
	}
	values, err = solve(t, matrix.NewSystem(a, c, d).Add(captured)).Bind(map[matrix.Var[string]]matrix.Rat{
		c: matrix.NewRat(3, 1), d: matrix.NewRat(2, 1),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 1, 1)
}

func TestExpressionBranchesAcrossTermCounts(t *testing.T) {
	a, b, c := matrix.Variable("a"), matrix.Variable("b"), matrix.Variable("c")
	base := b.Scale(2).Add(a.Scale(1))
	captured := base.Equal(5)
	child := base.Add(c.Scale(3))
	reduced := child.Add(b.Scale(-2))
	sibling := base.Add(a.Scale(4))
	for _, tt := range []struct {
		expression matrix.Expr[string]
		want       [3]int64
	}{
		{base, [3]int64{1, 2, 0}},
		{child, [3]int64{1, 2, 3}},
		{reduced, [3]int64{1, 0, 3}},
		{sibling, [3]int64{5, 2, 0}},
	} {
		for i, variable := range []matrix.Var[string]{a, b, c} {
			wantRat(t, tt.expression.Coefficient(variable), tt.want[i], 1)
		}
	}
	values, err := solve(t, matrix.NewSystem(a, b).Add(captured)).Bind(map[matrix.Var[string]]matrix.Rat{b: matrix.NewRat(2, 1)})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 1, 1)
}

func TestExpressionFractionalBranchesAndSnapshots(t *testing.T) {
	a, b := matrix.Variable("a"), matrix.Variable("b")
	half := a.ScaleRat(matrix.NewRat(1, 2))
	capturedHalf := half.Equal(1)
	pair := half.Add(b.Scale(2))
	reversed := b.Scale(2).Add(half)
	capturedPair := pair.Equal(5)
	whole := half.Add(a.ScaleRat(matrix.NewRat(1, 2)))
	wholePair := pair.Add(a.ScaleRat(matrix.NewRat(1, 2)))
	zero := half.Add(a.ScaleRat(matrix.NewRat(-1, 2)))
	for _, tt := range []struct {
		name                                   string
		expression                             matrix.Expr[string]
		aNumerator, aDenominator, bCoefficient int64
	}{
		{"original half", half, 1, 2, 0},
		{"original pair", pair, 1, 2, 2},
		{"reversed pair", reversed, 1, 2, 2},
		{"whole", whole, 1, 1, 0},
		{"whole pair", wholePair, 1, 1, 2},
		{"canceled", zero, 0, 1, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			wantRat(t, tt.expression.Coefficient(a), tt.aNumerator, tt.aDenominator)
			wantRat(t, tt.expression.Coefficient(b), tt.bCoefficient, 1)
		})
	}
	values, err := solve(t, matrix.NewSystem(a).Add(capturedHalf)).Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 2, 1)
	values, err = solve(t, matrix.NewSystem(a, b).Add(capturedPair)).Bind(map[matrix.Var[string]]matrix.Rat{
		b: matrix.NewRat(2, 1),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 2, 1)
	values, err = solve(t, matrix.NewSystem[string]().Add(zero.Equal(0))).Bind(nil)
	if err != nil || len(values) != 0 {
		t.Fatalf("canceled fractional expression = %v, %v", values, err)
	}
}

func TestSystemExtensionKeepsUniqueSnapshot(t *testing.T) {
	a, b := matrix.Variable("a"), matrix.Variable("b")
	system := matrix.NewSystem(b, a).
		Add(a.Scale(2).Add(b.Scale(1)).Equal(0)).
		Add(a.Scale(1).Add(b.Scale(2)).Equal(1))
	first := solve(t, system)
	second := solve(t, system.Add(a.Scale(4).Add(b.Scale(2)).Equal(0)))
	for _, solution := range []*matrix.Solution[string]{first, second} {
		if len(solution.Free()) != 0 {
			t.Fatal("redundant equation changed uniqueness")
		}
		values, err := solution.Bind(nil)
		if err != nil {
			t.Fatal(err)
		}
		wantRat(t, values[a], -1, 3)
		wantRat(t, values[b], 2, 3)
	}
}

func TestSystemOwnsLargerDeclarationSnapshot(t *testing.T) {
	a, b, c := matrix.Variable("a"), matrix.Variable("b"), matrix.Variable("c")
	variables := []matrix.Var[string]{b, a, c}
	system := matrix.NewSystem(variables...)
	variables[0] = matrix.Variable("changed")
	values, err := solve(t, system.Add(a.Scale(1).Equal(1)).Add(b.Scale(1).Equal(2)).Add(c.Scale(1).Equal(3))).Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 1, 1)
	wantRat(t, values[b], 2, 1)
	wantRat(t, values[c], 3, 1)
}

func TestSystemRepeatedSingleVariableConstraints(t *testing.T) {
	a := matrix.Variable("a")
	system := matrix.NewSystem(a).Add(a.Scale(1).Equal(3)).Add(a.Scale(2).Equal(6))
	first := solve(t, system)
	second := solve(t, system.Add(a.Scale(3).Equal(9)))
	for _, solution := range []*matrix.Solution[string]{first, second} {
		values, err := solution.Bind(nil)
		if err != nil {
			t.Fatal(err)
		}
		wantRat(t, values[a], 3, 1)
	}
}

func TestSystemRetainsConstantContradiction(t *testing.T) {
	var zero matrix.Expr[string]
	solution, err := matrix.NewSystem[string]().
		Add(zero.Equal(0)).Add(zero.Equal(1)).Add(zero.Equal(0)).Solve()
	if solution != nil || !errors.Is(err, matrix.ErrInconsistent) {
		t.Fatalf("contradictory constant equations = %v, %v", solution, err)
	}
}

func TestSystemErrors(t *testing.T) {
	a, b := matrix.Variable("a"), matrix.Variable("b")
	for _, tt := range []struct {
		name         string
		system       *matrix.System[string]
		inconsistent bool
		wantContext  string
	}{
		{"duplicate", matrix.NewSystem(a, a).Add(b.Scale(1).Equal(0)), false, "duplicate variable a"},
		{"undeclared", matrix.NewSystem(a).Add(b.Scale(1).Equal(0)), false, ""},
		{"inconsistent", matrix.NewSystem(a, b).
			Add(a.Scale(2).Add(b.Scale(1)).Equal(4)).
			Add(a.Scale(4).Add(b.Scale(2)).Equal(9)), true, ""},
		{"extra contradictory equation", matrix.NewSystem(a, b).
			Add(a.Scale(1).Equal(1)).Add(b.Scale(1).Equal(2)).
			Add(a.Scale(1).Equal(3)), true, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			solution, err := tt.system.Solve()
			if solution != nil || err == nil || errors.Is(err, matrix.ErrInconsistent) != tt.inconsistent {
				t.Fatalf("Solve = %v, %v; inconsistent = %v", solution, err, tt.inconsistent)
			}
			if tt.wantContext != "" && !strings.Contains(err.Error(), tt.wantContext) {
				t.Errorf("got %v; want %q context", err, tt.wantContext)
			}
		})
	}
	system := matrix.NewSystem(a).Add(b.Scale(1).Equal(0))
	_, first := system.Solve()
	system.Add(matrix.Variable("c").Scale(1).Equal(0))
	_, second := system.Solve()
	for _, err := range []error{first, second} {
		if err == nil || errors.Is(err, matrix.ErrInconsistent) ||
			!strings.Contains(err.Error(), "equation 1") || !strings.Contains(err.Error(), "variable b") {
			t.Fatalf("first equation error lost: %v", err)
		}
	}
}

func TestSystemValidationPrecedesReduction(t *testing.T) {
	a, b := matrix.Variable("a"), matrix.Variable("b")
	z, y := matrix.Variable("z"), matrix.Variable("y")
	var zero matrix.Expr[string]
	system := matrix.NewSystem(a, b).Add(zero.Equal(1)).Add(z.Scale(1).Equal(0))
	checkError := func(phase string) {
		t.Helper()
		solution, err := system.Solve()
		if solution != nil || err == nil || errors.Is(err, matrix.ErrInconsistent) ||
			!strings.Contains(err.Error(), "equation 2") || !strings.Contains(err.Error(), "variable z") {
			t.Fatalf("%s: got %v, %v; want first undeclared variable z in equation 2", phase, solution, err)
		}
	}
	checkError("before extension")
	system.Add(y.Scale(1).Equal(0))
	checkError("after extension")
	checkError("repeated solve")
}

func TestSolutionBindingErrors(t *testing.T) {
	a, b, unknown := matrix.Variable("a"), matrix.Variable("b"), matrix.Variable("unknown")
	solution := solve(t, matrix.NewSystem(a, b).Add(a.Scale(2).Add(b.Scale(1)).Equal(4)))
	for _, assignments := range []map[matrix.Var[string]]matrix.Rat{
		nil,
		{a: matrix.NewRat(1, 1)},
		{unknown: matrix.NewRat(1, 1)},
		{b: matrix.NewRat(1, 1), a: matrix.NewRat(1, 1)},
		{b: matrix.NewRat(1, 1), unknown: matrix.NewRat(1, 1)},
	} {
		values, err := solution.Bind(assignments)
		if values != nil || err == nil {
			t.Fatalf("Bind(%v) = %v, %v; want error without partial values", assignments, values, err)
		}
	}
	if _, err := solution.Expression(unknown); err == nil {
		t.Fatal("accepted undeclared variable")
	}
}

func TestSystemSnapshotsAndExpressionBranches(t *testing.T) {
	a, b, unknown := matrix.Variable("a"), matrix.Variable("b"), matrix.Variable("unknown")
	variables := []matrix.Var[string]{a, b}
	system := matrix.NewSystem(variables...)
	variables[0] = unknown
	var zero matrix.Expr[string]
	base := zero.Add(a.Scale(2)).Add(zero)
	captured := base.Equal(4)
	left, right := base.Add(b.Scale(1)), base.Add(b.Scale(3))
	solution := solve(t, system.Add(left.Equal(4)))
	sibling := solve(t, matrix.NewSystem(a, b).Add(right.Equal(4)))
	free := solution.Free()
	free[0] = unknown
	unique := solve(t, system.Add(b.Scale(1).Equal(2)))
	if len(unique.Free()) != 0 {
		t.Fatal("builder extension should fix b")
	}
	assignments := map[matrix.Var[string]]matrix.Rat{b: matrix.NewRat(1, 1)}
	values, err := solution.Bind(assignments)
	if err != nil {
		t.Fatal(err)
	}
	assignments[b] = matrix.NewRat(9, 1)
	wantRat(t, values[a], 3, 2)
	wantRat(t, values[b], 1, 1)
	values[a] = matrix.NewRat(99, 1)
	values, err = solution.Bind(map[matrix.Var[string]]matrix.Rat{b: matrix.NewRat(1, 1)})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 3, 2)
	values, err = sibling.Bind(map[matrix.Var[string]]matrix.Rat{b: matrix.NewRat(1, 1)})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 1, 2)
	wantRat(t, base.Coefficient(b), 0, 1)
	values, err = solve(t, matrix.NewSystem(a).Add(captured)).Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 2, 1)
}

func TestSystemTwoVariableEmptyEquation(t *testing.T) {
	a, b := matrix.Variable("a"), matrix.Variable("b")
	var empty matrix.Expr[string]
	for _, tt := range []struct {
		name    string
		rhs     int64
		wantErr error
	}{
		{"identity", 0, nil},
		{"contradiction", 1, matrix.ErrInconsistent},
	} {
		t.Run(tt.name, func(t *testing.T) {
			solution, err := matrix.NewSystem(b, a).
				Add(empty.Equal(tt.rhs)).
				Add(a.Scale(1).Equal(3)).
				Solve()
			if tt.wantErr != nil {
				if solution != nil || !errors.Is(err, tt.wantErr) {
					t.Fatalf("Solve = %v, %v; want nil, %v", solution, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(solution.Free(), []matrix.Var[string]{b}) {
				t.Fatal("expected b to be free")
			}
			values, err := solution.Bind(map[matrix.Var[string]]matrix.Rat{b: matrix.NewRat(2, 1)})
			if err != nil {
				t.Fatal(err)
			}
			wantRat(t, values[a], 3, 1)
			wantRat(t, values[b], 2, 1)
		})
	}
}

func TestSystemDegenerate(t *testing.T) {
	a, b := matrix.Variable("a"), matrix.Variable("b")
	var zero matrix.Expr[string]
	for _, system := range []*matrix.System[string]{
		{}, matrix.NewSystem[string]().Add(zero.Equal(0)),
		matrix.NewSystem[string]().Add(a.Scale(1).Add(a.Scale(-1)).Equal(0)),
	} {
		values, err := solve(t, system).Bind(nil)
		if err != nil || len(values) != 0 {
			t.Fatalf("empty system = %v, %v", values, err)
		}
	}
	if _, err := matrix.NewSystem[string]().Add(zero.Equal(1)).Solve(); !errors.Is(err, matrix.ErrInconsistent) {
		t.Fatalf("0=1: got %v", err)
	}
	allFree := solve(t, matrix.NewSystem(a, b))
	if !slices.Equal(allFree.Free(), []matrix.Var[string]{a, b}) {
		t.Fatal("equationless system lost declared variables")
	}
	values, err := allFree.Bind(map[matrix.Var[string]]matrix.Rat{a: matrix.NewRat(-1, 2), b: matrix.NewRat(7, 1)})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], -1, 2)
	wantRat(t, values[b], 7, 1)
	unused := solve(t, matrix.NewSystem(a, b).Add(a.Scale(1).Equal(3)))
	if !slices.Equal(unused.Free(), []matrix.Var[string]{b}) {
		t.Fatal("unused b must remain free")
	}
}
