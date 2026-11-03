package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type mapInstruction struct {
	steps int
	turn  byte
}

type monkeyMap struct {
	rows         []string
	width        int
	start        grid.Position2D[int]
	instructions []mapInstruction
}

func parseMonkeyMap(r io.Reader) (monkeyMap, error) {
	var board monkeyMap
	input, err := puzzleInput(r)
	if err != nil {
		return board, err
	}
	drawing, path, ok := strings.Cut(input, "\n\n")
	if !ok || len(drawing) == 0 || len(drawing) > 1_000_000 {
		return board, fmt.Errorf("map requires a drawing of at most one million characters and a path")
	}
	board.rows = strings.Split(drawing, "\n")
	for _, row := range board.rows {
		board.width = max(board.width, len(row))
		for _, cell := range row {
			if cell != '.' && cell != '#' && cell != ' ' {
				return board, fmt.Errorf("invalid map cell %q", cell)
			}
		}
	}
	board.start.X = strings.IndexByte(board.rows[0], '.')
	if board.start.X < 0 {
		return board, fmt.Errorf("map top row has no open starting tile")
	}
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return board, fmt.Errorf("empty map path")
	}
	for i := 0; i < len(path); {
		start := i
		for i < len(path) && path[i] >= '0' && path[i] <= '9' {
			i++
		}
		if start == i {
			return board, fmt.Errorf("expected forward count at path position %d", i+1)
		}
		steps, err := convert.ScanNumber[int]([]byte(path[start:i]))
		if err != nil {
			return board, err
		}
		instruction := mapInstruction{steps: steps}
		if i < len(path) {
			instruction.turn = path[i]
			if instruction.turn != 'L' && instruction.turn != 'R' || i+1 == len(path) {
				return board, fmt.Errorf("invalid turn at path position %d", i+1)
			}
			i++
		}
		board.instructions = append(board.instructions, instruction)
	}
	return board, nil
}

func (board monkeyMap) tile(position grid.Position2D[int]) byte {
	if position.Y < 0 || position.Y >= len(board.rows) || position.X < 0 || position.X >= len(board.rows[position.Y]) {
		return ' '
	}
	return board.rows[position.Y][position.X]
}

func mapDirection(facing int) grid.Position2D[int] {
	return [...]grid.Position2D[int]{{X: 1}, {Y: 1}, {X: -1}, {Y: -1}}[facing]
}

func (board monkeyMap) flatWrap(position grid.Position2D[int], facing int, work *int) (grid.Position2D[int], int, error) {
	direction := mapDirection(facing)
	for {
		*work--
		if *work < 0 {
			return position, facing, fmt.Errorf("map traversal exceeds ten million steps and void checks")
		}
		position.X = xmath.Modulo(position.X+direction.X, board.width)
		position.Y = xmath.Modulo(position.Y+direction.Y, len(board.rows))
		if board.tile(position) != ' ' {
			return position, facing, nil
		}
	}
}

func (board monkeyMap) password(wrap func(grid.Position2D[int], int, *int) (grid.Position2D[int], int, error)) (string, error) {
	position, facing := board.start, 0
	work := 10_000_000
	for _, instruction := range board.instructions {
		for range instruction.steps {
			work--
			if work < 0 {
				return "", fmt.Errorf("map traversal exceeds ten million steps and void checks")
			}
			next, nextFacing := position.Add(mapDirection(facing)), facing
			if board.tile(next) == ' ' {
				var err error
				next, nextFacing, err = wrap(position, facing, &work)
				if err != nil {
					return "", err
				}
			}
			if board.tile(next) == '#' {
				break
			}
			if board.tile(next) != '.' {
				return "", fmt.Errorf("map transition does not reach a tile")
			}
			position, facing = next, nextFacing
		}
		if instruction.turn == 'R' {
			facing = (facing + 1) % 4
		} else if instruction.turn == 'L' {
			facing = (facing + 3) % 4
		}
	}
	return strconv.Itoa(1000*(position.Y+1) + 4*(position.X+1) + facing), nil
}

type cubeFrame struct {
	right, down, normal grid.Position3D[int]
}

func negateCubeVector(vector grid.Position3D[int]) grid.Position3D[int] {
	return grid.Position3D[int]{X: -vector.X, Y: -vector.Y, Z: -vector.Z}
}

func (frame cubeFrame) direction(facing int) grid.Position3D[int] {
	return [...]grid.Position3D[int]{frame.right, frame.down, negateCubeVector(frame.right), negateCubeVector(frame.down)}[facing]
}

func (frame cubeFrame) fold(facing int) cubeFrame {
	switch facing {
	case 0:
		return cubeFrame{right: negateCubeVector(frame.normal), down: frame.down, normal: frame.right}
	case 1:
		return cubeFrame{right: frame.right, down: negateCubeVector(frame.normal), normal: frame.down}
	case 2:
		return cubeFrame{right: frame.normal, down: frame.down, normal: negateCubeVector(frame.right)}
	default:
		return cubeFrame{right: frame.right, down: frame.normal, normal: negateCubeVector(frame.down)}
	}
}

type cubeFace struct {
	origin grid.Position2D[int]
	frame  cubeFrame
}

type cubeNet struct {
	side  int
	faces []cubeFace
}

func foldCube(board monkeyMap) (cubeNet, error) {
	var cube cubeNet
	tiles := 0
	offset := board.start
	for y, row := range board.rows {
		for x, cell := range row {
			if cell != ' ' {
				tiles++
				offset.X, offset.Y = min(offset.X, x), min(offset.Y, y)
			}
		}
	}

	cube.side = 1
	for 6*cube.side*cube.side < tiles {
		cube.side++
	}
	if 6*cube.side*cube.side != tiles {
		return cube, fmt.Errorf("cube map requires six equal square faces")
	}

	byBlock := make(map[grid.Position2D[int]]int)
	var blocks []grid.Position2D[int]
	var counts []int
	for y, row := range board.rows {
		for x, cell := range row {
			if cell == ' ' {
				continue
			}
			block := grid.Position2D[int]{X: (x - offset.X) / cube.side, Y: (y - offset.Y) / cube.side}
			index, exists := byBlock[block]
			if !exists {
				index = len(blocks)
				byBlock[block] = index
				blocks = append(blocks, block)
				counts = append(counts, 0)
				cube.faces = append(cube.faces, cubeFace{origin: grid.Position2D[int]{X: offset.X + block.X*cube.side, Y: offset.Y + block.Y*cube.side}})
			}
			counts[index]++
		}
	}
	if len(blocks) != 6 {
		return cube, fmt.Errorf("cube tiles do not form six aligned square faces")
	}
	for _, count := range counts {
		if count != cube.side*cube.side {
			return cube, fmt.Errorf("cube face is incomplete")
		}
	}

	startBlock := grid.Position2D[int]{X: (board.start.X - offset.X) / cube.side, Y: (board.start.Y - offset.Y) / cube.side}
	start := byBlock[startBlock]
	cube.faces[start].frame = cubeFrame{right: grid.Position3D[int]{X: 1}, down: grid.Position3D[int]{Y: 1}, normal: grid.Position3D[int]{Z: 1}}
	assigned := [6]bool{}
	assigned[start] = true
	queue := []int{start}
	for head := 0; head < len(queue); head++ {
		index := queue[head]
		for facing := range 4 {
			block := blocks[index].Add(mapDirection(facing))
			next, exists := byBlock[block]
			if !exists {
				continue
			}
			frame := cube.faces[index].frame.fold(facing)
			if assigned[next] {
				if cube.faces[next].frame != frame {
					return cube, fmt.Errorf("cube net has inconsistent folds")
				}
				continue
			}
			assigned[next] = true
			cube.faces[next].frame = frame
			queue = append(queue, next)
		}
	}
	if len(queue) != 6 {
		return cube, fmt.Errorf("cube net is disconnected")
	}

	normals := make(map[grid.Position3D[int]]bool)
	for _, face := range cube.faces {
		if normals[face.frame.normal] {
			return cube, fmt.Errorf("cube net folds two faces onto the same side")
		}
		normals[face.frame.normal] = true
	}
	return cube, nil
}

func (cube cubeNet) wrap(position grid.Position2D[int], facing int, _ *int) (grid.Position2D[int], int, error) {
	var source cubeFace
	found := false
	for _, face := range cube.faces {
		if position.X >= face.origin.X && position.X < face.origin.X+cube.side && position.Y >= face.origin.Y && position.Y < face.origin.Y+cube.side {
			source, found = face, true
			break
		}
	}
	if !found {
		return position, facing, fmt.Errorf("position is not on a cube face")
	}
	for _, target := range cube.faces {
		if target.frame.normal != source.frame.direction(facing) {
			continue
		}
		// Crossing an edge turns forward motion toward the old face's inward normal.
		forward := negateCubeVector(source.frame.normal)
		for nextFacing := range 4 {
			if target.frame.direction(nextFacing) != forward {
				continue
			}
			parameter, axis := position.Y-source.origin.Y, source.frame.down
			if facing%2 == 1 {
				parameter, axis = position.X-source.origin.X, source.frame.right
			}
			nextAxis := target.frame.down
			if nextFacing%2 == 1 {
				nextAxis = target.frame.right
			}
			if axis == negateCubeVector(nextAxis) {
				parameter = cube.side - 1 - parameter
			} else if axis != nextAxis {
				return position, facing, fmt.Errorf("cube edges have incompatible axes")
			}
			local := [...]grid.Position2D[int]{
				{X: 0, Y: parameter},
				{X: parameter, Y: 0},
				{X: cube.side - 1, Y: parameter},
				{X: parameter, Y: cube.side - 1},
			}[nextFacing]
			return target.origin.Add(local), nextFacing, nil
		}
	}
	return position, facing, fmt.Errorf("cube edge has no destination")
}

func day22p01(r io.Reader) (string, error) {
	board, err := parseMonkeyMap(r)
	if err != nil {
		return "", err
	}
	return board.password(board.flatWrap)
}

func day22p02(r io.Reader) (string, error) {
	board, err := parseMonkeyMap(r)
	if err != nil {
		return "", err
	}
	cube, err := foldCube(board)
	if err != nil {
		return "", err
	}
	return board.password(cube.wrap)
}
