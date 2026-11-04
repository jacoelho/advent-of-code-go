// Package matrix builds and solves exact rational linear systems.
//
// Declare unknowns with Variable and NewSystem, compose immutable linear
// expressions, and submit complete equations with Add. Solve validates variable
// references in equation insertion order before reduction and identifies free
// variables; callers never handle pivot or column indexes.
// A free variable is an original unknown whose value parameterises a solution
// family, not an unknown marked free before solving.
//
// Rat is the canonical arithmetic representation. Generic Scale and Equal
// methods convert signed integer inputs exactly; ScaleRat and EqualRat accept
// calculated rational or arbitrary-precision coefficients. Bind produces
// rational values, without imposing integer feasibility, bounds or an objective.
//
// Systems run synchronously with no search or background work. Callers must
// bound variable counts, equation counts and input bit lengths. Add retains
// equations even after an invalid variable reference. Elimination uses dense
// storage and arbitrary-precision arithmetic whose cost grows with those
// dimensions and intermediate values. Puzzle-specific search budgets remain
// with the calling exercise.
package matrix
