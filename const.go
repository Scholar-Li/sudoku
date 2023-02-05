package sudoku

import (
	"strconv"
	"strings"
)

// 九宫格各格子的坐标值
var GridXY = [9][9]Coord{}

func init() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			GridXY[GetGrid(j, i)][GetGridIdx(j, i)] = Coord{
				X: j,
				Y: i,
			}
		}
	}
}

// 坐标轴，X 表示列，Y 表示行
type Coord struct {
	X, Y int
}

// 根据坐标轴，获取位于九宫格的第几个格子，
// 返回值范围 0-8
func GetGrid(x, y int) int {
	// y/3*3 != y，因为是先整除 3，取整数再乘以 3
	return y/3*3 + x/3
}

// 根据坐标轴，获取位于九宫格格子中的第几个下标，
// 返回值范围 0-8
func GetGridIdx(x, y int) int {
	return y%3*3 + x%3
}

// 数独 9*9 棋盘
type Chessboard [9][9]int

func (c Chessboard) String() string {
	res := make([]string, 9)

	for i := 0; i < 9; i++ {
		subRes := make([]string, 9)
		for j := 0; j < 9; j++ {
			subRes[j] = strconv.FormatInt(int64(c[i][j]), 10)
		}
		res[i] = strings.Join(subRes, " ")
	}
	return strings.Join(res, "\n")
}

// 检查九宫格，第 i 个格子是否已经满足
func (c Chessboard) CheckGrid(i int) (ok bool) {
	var r rightSudoku
	for _, coord := range GridXY[i] {
		r, ok = r.addNum(c[coord.Y][coord.X])
		if !ok {
			return false
		}
	}
	return true
}

// 检查第 i 列是否已经满足
func (c Chessboard) CheckLine(i int) (ok bool) {
	var r rightSudoku
	for _, row := range c {
		r, ok = r.addNum(row[i])
		if !ok {
			return false
		}
	}
	return true
}

// 检查第 i 行是否已经满足
func (c Chessboard) CheckRow(i int) bool {
	var r rightSudoku
	return r.addNums(c[i])
}

// 检测当前数独是否已经全部正确
func (c Chessboard) CheckRightSudoku() bool {
	// 检测行、列、九宫格已完成项
	for i := 0; i < 9; i++ {
		if !c.CheckGrid(i) {
			return false
		}
		if !c.CheckLine(i) {
			return false
		}
		if !c.CheckRow(i) {
			return false
		}
	}
	return true
}

// 检测数独各数字是否合法，即只存在 0-9 的值，0 表示空白格子，返回 false 表示非法
func (c Chessboard) ValidSudoku() bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			num := c[i][j]
			if num > 9 || num < 0 {
				return false
			}
		}
	}
	return true
}

// 检测正常九宫格的数字，即 1-9 数字都存在且不重复
type rightSudoku [10]int

// 检测 num 是否在范围 1-9 内
func (r rightSudoku) checkNum(num int) bool {
	if num > 9 || num < 1 {
		return false
	}
	return true
}

// nums 必为 9 个数字，且仅当为 1-9 不重复时，返回 true
func (r rightSudoku) addNums(nums [9]int) (ok bool) {
	for _, num := range nums {
		_, ok = r.addNum(num)
		if !ok {
			return false
		}
		r[num] = 1
	}

	return true
}

// 添加 num，如果仍然满足九宫格规则，则返回 true，否则立即返回 false
func (r rightSudoku) addNum(num int) (rightSudoku, bool) {
	if !r.checkNum(num) {
		return r, false
	}
	if r[num] != 0 {
		return r, false
	}
	r[num] = 1
	return r, true
}
