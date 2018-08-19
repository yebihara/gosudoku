package main

import (
	"fmt"
	"bytes"
	"strconv"
	"os"
	"bufio"
)

/*
 * Cell
 */
type Cell struct {
	value int      // cell value
	given bool     // true if the value is given by the problem
	row *Group     // the row this cell belongs to
	col *Group     // the column this cell belongs to
	block *Group   // the block this cell belongs to
	next *Cell     // reference to the next cell (left to right, then top to bottom)
	sudoku *Sudoku // reference to the entire set of cells
}

func (c *Cell) CheckViolation() bool {
	return c.row.CheckViolation() && c.col.CheckViolation() && c.block.CheckViolation()
}

/*
 *  Look for the answer for the cell based on backtracking algorithm
 */
func (cell *Cell) Test1to9(debug bool) bool {
	if cell == nil {
		// stop recursion at the last cell
		return true
	}

	if cell.given {
		return cell.next.Test1to9(debug)
	} else {
		// test through 1 to 9
		for value := 1; value <= 9; value++ {
			cell.value = value

			if cell.CheckViolation() {
				if debug {
					fmt.Println(cell.sudoku)
				}

				if cell.next.Test1to9(debug) {
					return true
				}
			}
		}
	}

	cell.value = 0
	return false
}

func (c Cell) String() string {
	if c.value == 0 {
		return "_"
	} else {
		return strconv.Itoa(c.value)
	}
}

/*
 * Group
 */
type Group [9]*Cell

func (g Group) CheckViolation() bool {
	checker := make(map[int]bool)
	for _, cell := range g {
		if cell.value > 0 {
			exists := checker[cell.value]
			if exists {
				return false
			} else {
				checker[cell.value] = true
			}
		}
	}

	return true
}

/*
 * Sudoku
 */
type Sudoku struct {
	cells [81]*Cell
	rows [9]*Group
	cols [9]*Group
	blocks [9]*Group
}

func (s Sudoku) String() string {
	var buf bytes.Buffer

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			buf.WriteString(s.cells[i * 9 + j].String())
			buf.WriteString(" ")
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func (s Sudoku) Solve(debug bool) bool {
	return s.cells[0].Test1to9(debug)
}

/*
 * Initialize
 */
func Initialize(problem [81]int) *Sudoku {
	sudoku := Sudoku{}

	for i := 0; i < 9; i++ {
		sudoku.rows[i] = &Group{}
		sudoku.cols[i] = &Group{}
		sudoku.blocks[i] = &Group{}
	}

	var prev_cell *Cell

	for idx, value := range problem {
		row_no := idx / 9
		col_no := idx % 9
		blk_no := row_no / 3 * 3 + col_no / 3

		cell := Cell{value: value, given: value != 0, row: sudoku.rows[row_no], col: sudoku.cols[col_no], block: sudoku.blocks[blk_no], sudoku: &sudoku}

		sudoku.cells[idx] = &cell
		sudoku.rows[row_no][col_no] = &cell
		sudoku.cols[col_no][row_no] = &cell
		sudoku.blocks[blk_no][row_no % 3 * 3 + col_no % 3] = &cell

		if prev_cell != nil {
			// not set to the first cell
			prev_cell.next = &cell
		}
		prev_cell = &cell
	}

	return &sudoku
}

/*
Sample Data

_7_2_____
3_8__6_7_
__9_____5
___72____
25______8
____9_3__
__3_1___6
_4_5___3_
___6_781_
*/
func ReadProblem() [81]int {
	var problem [81]int

	stdin := bufio.NewScanner(os.Stdin)

	for i := 0; i < 81; {
		stdin.Scan()
		input := stdin.Text()

		for j := 0; j < len(input) && i < 81; j++ {
			b := input[j]

			// parse '1' to '9' as an integer, otherwise 0
			if b >= 0x31 && b <= 0x39 {
				problem[i] = int(b) - '0'
			} else {
				problem[i] = 0
			}

			i++
		}
	}

	return problem
}

/*
 * main
 */
func main() {
	problem := ReadProblem()

	sudoku := Initialize(problem)
	fmt.Println(sudoku)

	success := sudoku.Solve(false)

	fmt.Println("\nSuccess? => ", success, "\n")

	if success {
		fmt.Println(sudoku)
	}
}

