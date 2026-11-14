package linear_test

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/pkg/linear"
)

func ExampleSystem() {
	a, b := linear.Variable("a"), linear.Variable("b")
	solution, err := linear.NewSystem(a, b).
		Add(a.Scale(2).Add(b.Scale(1)).Equate(4)).
		Solve()
	if err != nil {
		fmt.Println(err)
		return
	}
	values, err := solution.Bind(map[linear.Var[string]]linear.Rat{
		b: linear.NewRat(2, 1),
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

func solve[K linear.Key](t *testing.T, system *linear.System[K]) *linear.Solution[K] {
	t.Helper()
	solution, err := system.Solve()
	if err != nil {
		t.Fatal(err)
	}
	return solution
}

func wantRat(t *testing.T, got linear.Rat, numerator, denominator int64) {
	t.Helper()
	if got.Cmp(linear.NewRat(numerator, denominator)) != 0 {
		t.Fatalf("got %s/%s; want %d/%d", got.Numerator(), got.Denominator(), numerator, denominator)
	}
}

func TestSystemUnique(t *testing.T) {
	type component uint8
	type coefficient int64
	x, y := linear.Variable(component(0)), linear.Variable(component(1))
	solution := solve(t, linear.NewSystem(x, y).
		Add(y.Scale(coefficient(2)).Equate(int64(8))).
		Add(x.Scale(3).Add(y.Scale(1)).Equate(7)))
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
	*solution = linear.Solution[component]{}
	wantRat(t, captured.Constant(), 1, 1)
}

func TestSystemRationalFamily(t *testing.T) {
	a, b := linear.Variable("a"), linear.Variable("b")
	solution := solve(t, linear.NewSystem(a, b).
		Add(a.Scale(2).Add(b.Scale(1)).Equate(4)).
		Add(a.Scale(4).Add(b.Scale(2)).Equate(8)))
	if !slices.Equal(solution.Free(), []linear.Var[string]{b}) {
		t.Fatal("expected b to be free")
	}
	expression, err := solution.Expression(linear.Variable("a"))
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, expression.Constant(), 2, 1)
	wantRat(t, expression.Coefficient(b), -1, 2)
	wantRat(t, expression.Coefficient(a), 0, 1)
	for _, binding := range []struct{ b, numerator, denominator int64 }{{1, 3, 2}, {2, 1, 1}} {
		values, err := solution.Bind(map[linear.Var[string]]linear.Rat{b: linear.NewRat(binding.b, 1)})
		if err != nil {
			t.Fatal(err)
		}
		wantRat(t, values[a], binding.numerator, binding.denominator)
		wantRat(t, values[b], binding.b, 1)
	}
	bound := solve(t, linear.NewSystem(b).Add(expression.Equate(1)))
	values, err := bound.Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[b], 2, 1)

	exact := solve(t, linear.NewSystem(a).
		Add(a.ScaleRat(linear.NewRat(1, 2)).EquateRat(linear.NewRat(3, 2))))
	values, err = exact.Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 3, 1)
}

func TestAffineExpressionOperationsAndSnapshots(t *testing.T) {
	x, z, parameter := linear.Variable("x"), linear.Variable("z"), linear.Variable("t")
	solution := solve(t, linear.NewSystem(x, z, parameter).
		Add(x.Scale(1).Add(parameter.Scale(1)).Equate(2)).
		Add(x.Scale(1).Add(z.Scale(1)).Equate(0)))
	expressionX, err := solution.Expression(x)
	if err != nil {
		t.Fatal(err)
	}
	expressionZ, err := solution.Expression(z)
	if err != nil {
		t.Fatal(err)
	}
	captured := expressionX.Equate(1)
	doubled := expressionX.Add(expressionX)
	canceled := expressionX.Add(expressionZ)
	for _, tt := range []struct {
		name                  string
		expression            linear.Expr[string]
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
	values, err := solve(t, linear.NewSystem(parameter).Add(captured)).Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[parameter], 1, 1)
	values, err = solve(t, linear.NewSystem[string]().Add(canceled.Equate(0))).Bind(nil)
	if err != nil || len(values) != 0 {
		t.Fatalf("canceled affine expression = %v, %v", values, err)
	}
}

func TestSystemTwoVariableExactSolutions(t *testing.T) {
	x, y := linear.Variable("x"), linear.Variable("y")
	for _, tt := range []struct {
		name   string
		first  linear.Equation[string]
		second linear.Equation[string]
		wantX  linear.Rat
		wantY  linear.Rat
	}{
		{
			"fractional signed solution",
			x.Scale(2).Add(y.Scale(1)).Equate(0),
			x.Scale(1).Add(y.Scale(2)).Equate(1),
			linear.NewRat(-1, 3), linear.NewRat(2, 3),
		},
		{
			"fractional coefficient",
			x.ScaleRat(linear.NewRat(1, 2)).Add(y.Scale(1)).Equate(2),
			x.Scale(1).Add(y.Scale(3)).Equate(5),
			linear.NewRat(2, 1), linear.NewRat(1, 1),
		},
		{
			"fractional right hand side",
			x.Scale(2).Add(y.Scale(1)).EquateRat(linear.NewRat(1, 2)),
			x.Scale(1).Add(y.Scale(2)).Equate(1),
			linear.NewRat(0, 1), linear.NewRat(1, 2),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			solution := solve(t, linear.NewSystem(x, y).Add(tt.first).Add(tt.second))
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
	x, y := linear.Variable("x"), linear.Variable("y")
	large := new(big.Int).Lsh(big.NewInt(1), 80)
	coefficient := linear.NewInteger(large)
	system := linear.NewSystem(x, y).
		Add(x.ScaleRat(coefficient).EquateRat(linear.NewInteger(new(big.Int).Mul(large, big.NewInt(2))))).
		Add(y.ScaleRat(coefficient).EquateRat(linear.NewInteger(new(big.Int).Mul(large, big.NewInt(3)))))
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
	x, y := linear.Variable("x"), linear.Variable("y")
	for _, tt := range []struct {
		name          string
		first, second linear.Equation[string]
		wantX, wantY  string
	}{
		{"determinant subtraction", x.Scale(int64(math.MaxInt64)).Add(y.Scale(1)).Equate(1),
			x.Scale(int64(math.MinInt64)).Add(y.Scale(1)).Equate(1), "0", "1"},
		{"first numerator product", x.Scale(1).Add(y.Scale(2)).Equate(0),
			y.Scale(1).Equate(int64(math.MaxInt64)), "-18446744073709551614", "9223372036854775807"},
		{"second numerator product", x.Scale(1).Equate(int64(math.MaxInt64)),
			x.Scale(2).Add(y.Scale(1)).Equate(0), "9223372036854775807", "-18446744073709551614"},
		{"numerator subtraction", x.Scale(1).Add(y.Scale(-1)).Equate(int64(math.MaxInt64)),
			y.Scale(1).Equate(1), "9223372036854775808", "1"},
		{"minimum divided by negative one", x.Scale(-1).Equate(int64(math.MinInt64)),
			y.Scale(1).Equate(0), "9223372036854775808", "0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			values, err := solve(t, linear.NewSystem(x, y).Add(tt.first).Add(tt.second)).Bind(nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []struct {
				variable linear.Var[string]
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
	a, b, c, d := linear.Variable("a"), linear.Variable("b"), linear.Variable("c"), linear.Variable("d")
	system := linear.NewSystem(c, a, d, b).
		Add(d.Scale(3).Add(a.Scale(1)).Add(b.Scale(2)).Equate(7)).
		Add(c.Scale(1).Add(d.Scale(-4)).Equate(5))
	for _, solution := range []*linear.Solution[string]{solve(t, system), solve(t, system)} {
		if !slices.Equal(solution.Free(), []linear.Var[string]{d, b}) {
			t.Fatal("expected free variables d, b in declaration order")
		}
		expression, err := solution.Expression(a)
		if err != nil {
			t.Fatal(err)
		}
		wantRat(t, expression.Constant(), 7, 1)
		wantRat(t, expression.Coefficient(b), -2, 1)
		wantRat(t, expression.Coefficient(d), -3, 1)
		values, err := solution.Bind(map[linear.Var[string]]linear.Rat{b: linear.NewRat(1, 1), d: linear.NewRat(2, 1)})
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
	a, b, c, d := linear.Variable("a"), linear.Variable("b"), linear.Variable("c"), linear.Variable("d")
	left := d.Scale(4).Add(b.Scale(2)).Add(a.Scale(1))
	right := c.Scale(3).Add(b.Scale(-2)).Add(a.Scale(5))
	merged := left.Add(right)
	for _, tt := range []struct {
		expression linear.Expr[string]
		want       [4]int64
	}{
		{left, [4]int64{1, 2, 0, 4}},
		{right, [4]int64{5, -2, 3, 0}},
		{merged, [4]int64{6, 0, 3, 4}},
	} {
		for i, variable := range []linear.Var[string]{a, b, c, d} {
			wantRat(t, tt.expression.Coefficient(variable), tt.want[i], 1)
		}
	}
	captured := merged.Equate(23)
	canceled := merged.Add(a.Scale(-6)).Add(c.Scale(-3)).Add(d.Scale(-4))
	values, err := solve(t, linear.NewSystem[string]().Add(canceled.Equate(0))).Bind(nil)
	if err != nil || len(values) != 0 {
		t.Fatalf("canceled expression = %v, %v", values, err)
	}
	values, err = solve(t, linear.NewSystem(a, c, d).Add(captured)).Bind(map[linear.Var[string]]linear.Rat{
		c: linear.NewRat(3, 1), d: linear.NewRat(2, 1),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 1, 1)
}

func TestExpressionBranchesAcrossTermCounts(t *testing.T) {
	a, b, c := linear.Variable("a"), linear.Variable("b"), linear.Variable("c")
	base := b.Scale(2).Add(a.Scale(1))
	captured := base.Equate(5)
	child := base.Add(c.Scale(3))
	reduced := child.Add(b.Scale(-2))
	sibling := base.Add(a.Scale(4))
	for _, tt := range []struct {
		expression linear.Expr[string]
		want       [3]int64
	}{
		{base, [3]int64{1, 2, 0}},
		{child, [3]int64{1, 2, 3}},
		{reduced, [3]int64{1, 0, 3}},
		{sibling, [3]int64{5, 2, 0}},
	} {
		for i, variable := range []linear.Var[string]{a, b, c} {
			wantRat(t, tt.expression.Coefficient(variable), tt.want[i], 1)
		}
	}
	values, err := solve(t, linear.NewSystem(a, b).Add(captured)).Bind(map[linear.Var[string]]linear.Rat{b: linear.NewRat(2, 1)})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 1, 1)
}

func TestExpressionFractionalBranchesAndSnapshots(t *testing.T) {
	a, b := linear.Variable("a"), linear.Variable("b")
	half := a.ScaleRat(linear.NewRat(1, 2))
	capturedHalf := half.Equate(1)
	pair := half.Add(b.Scale(2))
	reversed := b.Scale(2).Add(half)
	capturedPair := pair.Equate(5)
	whole := half.Add(a.ScaleRat(linear.NewRat(1, 2)))
	wholePair := pair.Add(a.ScaleRat(linear.NewRat(1, 2)))
	zero := half.Add(a.ScaleRat(linear.NewRat(-1, 2)))
	for _, tt := range []struct {
		name                                   string
		expression                             linear.Expr[string]
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
	values, err := solve(t, linear.NewSystem(a).Add(capturedHalf)).Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 2, 1)
	values, err = solve(t, linear.NewSystem(a, b).Add(capturedPair)).Bind(map[linear.Var[string]]linear.Rat{
		b: linear.NewRat(2, 1),
	})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 2, 1)
	values, err = solve(t, linear.NewSystem[string]().Add(zero.Equate(0))).Bind(nil)
	if err != nil || len(values) != 0 {
		t.Fatalf("canceled fractional expression = %v, %v", values, err)
	}
}

func TestSystemExtensionKeepsUniqueSnapshot(t *testing.T) {
	a, b := linear.Variable("a"), linear.Variable("b")
	system := linear.NewSystem(b, a).
		Add(a.Scale(2).Add(b.Scale(1)).Equate(0)).
		Add(a.Scale(1).Add(b.Scale(2)).Equate(1))
	first := solve(t, system)
	second := solve(t, system.Add(a.Scale(4).Add(b.Scale(2)).Equate(0)))
	for _, solution := range []*linear.Solution[string]{first, second} {
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
	a, b, c := linear.Variable("a"), linear.Variable("b"), linear.Variable("c")
	variables := []linear.Var[string]{b, a, c}
	system := linear.NewSystem(variables...)
	variables[0] = linear.Variable("changed")
	values, err := solve(t, system.Add(a.Scale(1).Equate(1)).Add(b.Scale(1).Equate(2)).Add(c.Scale(1).Equate(3))).Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 1, 1)
	wantRat(t, values[b], 2, 1)
	wantRat(t, values[c], 3, 1)
}

func TestSystemRepeatedSingleVariableConstraints(t *testing.T) {
	a := linear.Variable("a")
	system := linear.NewSystem(a).Add(a.Scale(1).Equate(3)).Add(a.Scale(2).Equate(6))
	first := solve(t, system)
	second := solve(t, system.Add(a.Scale(3).Equate(9)))
	for _, solution := range []*linear.Solution[string]{first, second} {
		values, err := solution.Bind(nil)
		if err != nil {
			t.Fatal(err)
		}
		wantRat(t, values[a], 3, 1)
	}
}

func TestSystemRetainsConstantContradiction(t *testing.T) {
	var zero linear.Expr[string]
	solution, err := linear.NewSystem[string]().
		Add(zero.Equate(0)).Add(zero.Equate(1)).Add(zero.Equate(0)).Solve()
	if solution != nil || !errors.Is(err, linear.ErrInconsistent) {
		t.Fatalf("contradictory constant equations = %v, %v", solution, err)
	}
}

func TestSystemErrors(t *testing.T) {
	a, b := linear.Variable("a"), linear.Variable("b")
	for _, tt := range []struct {
		name         string
		system       *linear.System[string]
		inconsistent bool
		wantContext  string
	}{
		{"duplicate", linear.NewSystem(a, a).Add(b.Scale(1).Equate(0)), false, "duplicate variable a"},
		{"undeclared", linear.NewSystem(a).Add(b.Scale(1).Equate(0)), false, ""},
		{"inconsistent", linear.NewSystem(a, b).
			Add(a.Scale(2).Add(b.Scale(1)).Equate(4)).
			Add(a.Scale(4).Add(b.Scale(2)).Equate(9)), true, ""},
		{"extra contradictory equation", linear.NewSystem(a, b).
			Add(a.Scale(1).Equate(1)).Add(b.Scale(1).Equate(2)).
			Add(a.Scale(1).Equate(3)), true, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			solution, err := tt.system.Solve()
			if solution != nil || err == nil || errors.Is(err, linear.ErrInconsistent) != tt.inconsistent {
				t.Fatalf("Solve = %v, %v; inconsistent = %v", solution, err, tt.inconsistent)
			}
			if tt.wantContext != "" && !strings.Contains(err.Error(), tt.wantContext) {
				t.Errorf("got %v; want %q context", err, tt.wantContext)
			}
		})
	}
	system := linear.NewSystem(a).Add(b.Scale(1).Equate(0))
	_, first := system.Solve()
	system.Add(linear.Variable("c").Scale(1).Equate(0))
	_, second := system.Solve()
	for _, err := range []error{first, second} {
		if err == nil || errors.Is(err, linear.ErrInconsistent) ||
			!strings.Contains(err.Error(), "equation 1") || !strings.Contains(err.Error(), "variable b") {
			t.Fatalf("first equation error lost: %v", err)
		}
	}
}

func TestSystemValidationPrecedesReduction(t *testing.T) {
	a, b := linear.Variable("a"), linear.Variable("b")
	z, y := linear.Variable("z"), linear.Variable("y")
	var zero linear.Expr[string]
	system := linear.NewSystem(a, b).Add(zero.Equate(1)).Add(z.Scale(1).Equate(0))
	checkError := func(phase string) {
		t.Helper()
		solution, err := system.Solve()
		if solution != nil || err == nil || errors.Is(err, linear.ErrInconsistent) ||
			!strings.Contains(err.Error(), "equation 2") || !strings.Contains(err.Error(), "variable z") {
			t.Fatalf("%s: got %v, %v; want first undeclared variable z in equation 2", phase, solution, err)
		}
	}
	checkError("before extension")
	system.Add(y.Scale(1).Equate(0))
	checkError("after extension")
	checkError("repeated solve")
}

func TestSolutionBindingErrors(t *testing.T) {
	a, b, unknown := linear.Variable("a"), linear.Variable("b"), linear.Variable("unknown")
	solution := solve(t, linear.NewSystem(a, b).Add(a.Scale(2).Add(b.Scale(1)).Equate(4)))
	for _, assignments := range []map[linear.Var[string]]linear.Rat{
		nil,
		{a: linear.NewRat(1, 1)},
		{unknown: linear.NewRat(1, 1)},
		{b: linear.NewRat(1, 1), a: linear.NewRat(1, 1)},
		{b: linear.NewRat(1, 1), unknown: linear.NewRat(1, 1)},
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
	a, b, unknown := linear.Variable("a"), linear.Variable("b"), linear.Variable("unknown")
	variables := []linear.Var[string]{a, b}
	system := linear.NewSystem(variables...)
	variables[0] = unknown
	var zero linear.Expr[string]
	base := zero.Add(a.Scale(2)).Add(zero)
	captured := base.Equate(4)
	left, right := base.Add(b.Scale(1)), base.Add(b.Scale(3))
	solution := solve(t, system.Add(left.Equate(4)))
	sibling := solve(t, linear.NewSystem(a, b).Add(right.Equate(4)))
	free := solution.Free()
	free[0] = unknown
	unique := solve(t, system.Add(b.Scale(1).Equate(2)))
	if len(unique.Free()) != 0 {
		t.Fatal("builder extension should fix b")
	}
	assignments := map[linear.Var[string]]linear.Rat{b: linear.NewRat(1, 1)}
	values, err := solution.Bind(assignments)
	if err != nil {
		t.Fatal(err)
	}
	assignments[b] = linear.NewRat(9, 1)
	wantRat(t, values[a], 3, 2)
	wantRat(t, values[b], 1, 1)
	values[a] = linear.NewRat(99, 1)
	values, err = solution.Bind(map[linear.Var[string]]linear.Rat{b: linear.NewRat(1, 1)})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 3, 2)
	values, err = sibling.Bind(map[linear.Var[string]]linear.Rat{b: linear.NewRat(1, 1)})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 1, 2)
	wantRat(t, base.Coefficient(b), 0, 1)
	values, err = solve(t, linear.NewSystem(a).Add(captured)).Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], 2, 1)
}

func TestSystemTwoVariableEmptyEquation(t *testing.T) {
	a, b := linear.Variable("a"), linear.Variable("b")
	var empty linear.Expr[string]
	for _, tt := range []struct {
		name    string
		rhs     int64
		wantErr error
	}{
		{"identity", 0, nil},
		{"contradiction", 1, linear.ErrInconsistent},
	} {
		t.Run(tt.name, func(t *testing.T) {
			solution, err := linear.NewSystem(b, a).
				Add(empty.Equate(tt.rhs)).
				Add(a.Scale(1).Equate(3)).
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
			if !slices.Equal(solution.Free(), []linear.Var[string]{b}) {
				t.Fatal("expected b to be free")
			}
			values, err := solution.Bind(map[linear.Var[string]]linear.Rat{b: linear.NewRat(2, 1)})
			if err != nil {
				t.Fatal(err)
			}
			wantRat(t, values[a], 3, 1)
			wantRat(t, values[b], 2, 1)
		})
	}
}

func TestSystemDegenerate(t *testing.T) {
	a, b := linear.Variable("a"), linear.Variable("b")
	var zero linear.Expr[string]
	for _, system := range []*linear.System[string]{
		{}, linear.NewSystem[string]().Add(zero.Equate(0)),
		linear.NewSystem[string]().Add(a.Scale(1).Add(a.Scale(-1)).Equate(0)),
	} {
		values, err := solve(t, system).Bind(nil)
		if err != nil || len(values) != 0 {
			t.Fatalf("empty system = %v, %v", values, err)
		}
	}
	if _, err := linear.NewSystem[string]().Add(zero.Equate(1)).Solve(); !errors.Is(err, linear.ErrInconsistent) {
		t.Fatalf("0=1: got %v", err)
	}
	allFree := solve(t, linear.NewSystem(a, b))
	if !slices.Equal(allFree.Free(), []linear.Var[string]{a, b}) {
		t.Fatal("equationless system lost declared variables")
	}
	values, err := allFree.Bind(map[linear.Var[string]]linear.Rat{a: linear.NewRat(-1, 2), b: linear.NewRat(7, 1)})
	if err != nil {
		t.Fatal(err)
	}
	wantRat(t, values[a], -1, 2)
	wantRat(t, values[b], 7, 1)
	unused := solve(t, linear.NewSystem(a, b).Add(a.Scale(1).Equate(3)))
	if !slices.Equal(unused.Free(), []linear.Var[string]{b}) {
		t.Fatal("unused b must remain free")
	}
}
