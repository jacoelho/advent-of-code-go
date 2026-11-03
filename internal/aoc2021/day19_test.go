package aoc2021

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func Test_day19p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(scannerExample), Want: "79"},
		{Input: aoc.FileInput(t, 2021, 19), Want: "320"},
	}

	aoc.AOCTest(t, day19p01, tests)
}

func Test_day19p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(scannerExample), Want: "3621"},
		{Input: aoc.FileInput(t, 2021, 19), Want: "9655"},
	}

	aoc.AOCTest(t, day19p02, tests)
}

func TestScannerTranslation(t *testing.T) {
	var input strings.Builder
	input.WriteString("--- scanner 0 ---\n")
	for i := range 12 {
		fmt.Fprintf(&input, "%d,%d,%d\n", i, i*i, i*i*i)
	}
	input.WriteString("\n--- scanner 1 ---\n")
	for i := range 12 {
		fmt.Fprintf(&input, "%d,%d,%d\n", i-100, i*i+50, i*i*i-25)
	}
	count, distance, err := alignScanners(strings.NewReader(input.String()))
	if err != nil {
		t.Fatal(err)
	}
	if count != 12 || distance != 175 {
		t.Fatalf("alignment = %d beacons, distance %d; want 12, 175", count, distance)
	}
}

func TestScannerRotations(t *testing.T) {
	seen := make(map[grid.Position3D[int64]]bool)
	for _, rotation := range scannerRotations() {
		point := rotateBeacon(grid.NewPosition3D[int64](1, 2, 3), rotation)
		if seen[point] {
			t.Fatalf("duplicate orientation %v", point)
		}
		seen[point] = true
		x := rotateBeacon(grid.NewPosition3D[int64](1, 0, 0), rotation)
		y := rotateBeacon(grid.NewPosition3D[int64](0, 1, 0), rotation)
		z := rotateBeacon(grid.NewPosition3D[int64](0, 0, 1), rotation)
		cross := grid.NewPosition3D(x.Y*y.Z-x.Z*y.Y, x.Z*y.X-x.X*y.Z, x.X*y.Y-x.Y*y.X)
		if cross != z || x.X*x.X+x.Y*x.Y+x.Z*x.Z != 1 || y.X*y.X+y.Y*y.Y+y.Z*y.Z != 1 {
			t.Fatalf("orientation %v is not a proper rotation", rotation)
		}
	}
	if len(seen) != 24 {
		t.Fatalf("got %d orientations, want 24", len(seen))
	}
}

func TestScannerReflection(t *testing.T) {
	var input strings.Builder
	for scanner, sign := range []int{1, -1} {
		fmt.Fprintf(&input, "--- scanner %d ---\n", scanner)
		for i := range 12 {
			fmt.Fprintf(&input, "%d,%d,%d\n", sign*i, i*i, i*i*i)
		}
		input.WriteByte('\n')
	}
	if _, err := day19p01(strings.NewReader(input.String())); err == nil {
		t.Fatal("reflected scanners must not align")
	}
}

func TestScannerErrors(t *testing.T) {
	for _, input := range []string{
		"--- scanner 0 ---\n0,0,0\n0,0,0",
		"--- scanner 0 ---\n0,0,0\n\n--- scanner 1 ---\n1,1,1",
		"--- scanner 0 ---\n0,0,0\n\n--- scanner 0 ---\n1,1,1",
		"--- scanner 0 ---\n0,1",
		"--- scanner 0 ---",
	} {
		if _, err := day19p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

const scannerExample = `--- scanner 0 ---
404,-588,-901
528,-643,409
-838,591,734
390,-675,-793
-537,-823,-458
-485,-357,347
-345,-311,381
-661,-816,-575
-876,649,763
-618,-824,-621
553,345,-567
474,580,667
-447,-329,318
-584,868,-557
544,-627,-890
564,392,-477
455,729,728
-892,524,684
-689,845,-530
423,-701,434
7,-33,-71
630,319,-379
443,580,662
-789,900,-551
459,-707,401

--- scanner 1 ---
686,422,578
605,423,415
515,917,-361
-336,658,858
95,138,22
-476,619,847
-340,-569,-846
567,-361,727
-460,603,-452
669,-402,600
729,430,532
-500,-761,534
-322,571,750
-466,-666,-811
-429,-592,574
-355,545,-477
703,-491,-529
-328,-685,520
413,935,-424
-391,539,-444
586,-435,557
-364,-763,-893
807,-499,-711
755,-354,-619
553,889,-390

--- scanner 2 ---
649,640,665
682,-795,504
-784,533,-524
-644,584,-595
-588,-843,648
-30,6,44
-674,560,763
500,723,-460
609,671,-379
-555,-800,653
-675,-892,-343
697,-426,-610
578,704,681
493,664,-388
-671,-858,530
-667,343,800
571,-461,-707
-138,-166,112
-889,563,-600
646,-828,498
640,759,510
-630,509,768
-681,-892,-333
673,-379,-804
-742,-814,-386
577,-820,562

--- scanner 3 ---
-589,542,597
605,-692,669
-500,565,-823
-660,373,557
-458,-679,-417
-488,449,543
-626,468,-788
338,-750,-386
528,-832,-391
562,-778,733
-938,-730,414
543,643,-506
-524,371,-870
407,773,750
-104,29,83
378,-903,-323
-778,-728,485
426,699,580
-438,-605,-362
-469,-447,-387
509,732,623
647,635,-688
-868,-804,481
614,-800,639
595,780,-596

--- scanner 4 ---
727,592,562
-293,-554,779
441,611,-461
-714,465,-776
-743,427,-804
-660,-479,-426
832,-632,460
927,-485,-438
408,393,-506
466,436,-512
110,16,151
-258,-428,682
-393,719,612
-211,-452,876
808,-476,-593
-575,615,604
-485,667,467
-680,325,-822
-627,-443,-432
872,-547,-609
833,512,582
807,604,487
839,-516,451
891,-625,532
-652,-548,-490
30,-46,-14`
