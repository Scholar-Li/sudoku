package sudoku

// 获取数独/九宫格的解，如果返回 false 表示，数独无解
func Sudoku(chessboard Chessboard) (result Chessboard, ok bool) {
	// 初始化校验九宫格各数值
	if !chessboard.ValidSudoku() {
		return chessboard, false
	}

	// 回朔法求解
	chessboard, ok = backtracking(chessboard)
	if !ok {
		return chessboard, false
	}

	return chessboard, chessboard.CheckRightSudoku()
}

// 通过回溯法获取数独/九宫格的解：
//	1. 利用排除法，获取到各空白格有效填入值；
//	2. 通过不断枚举最优空白格的有效填入值，直至获取到正确解决
func backtracking(chessboard Chessboard) (Chessboard, bool) {
	// 获取九宫格最优解，若无解，则回溯
	chessboard, availableData, coord, ok := GetChessboardAvailableValues(chessboard)
	if !ok {
		return chessboard, false
	}

	// 若最优空白格存在值，即不存在空白格，代表当前九宫格都已赋值，已经得到解
	if chessboard[coord.Y][coord.X] != 0 {
		return chessboard, true
	}

	// 枚举最优空白格有效填入值，获取到正确解后立即返回
	for _, num := range availableData[coord.Y][coord.X] {
		newChessboard := chessboard
		newChessboard[coord.Y][coord.X] = num
		newChessboard, ok = backtracking(newChessboard)
		if ok {
			return newChessboard, true
		}
	}

	// 获取不到当前九宫格的解，回溯
	return chessboard, false
}

// 通过排除法获取九宫格各格子可填值列表：
//	res，最新的九宫格数值；
//	availableData，各坐标对应的有效填入值列表；
//	coord，可填值范围最小的坐标，即最优解；
//	ok，是否返回结果正常，若返回 false 表示给的九宫格已陷入无解中。
func GetChessboardAvailableValues(chessboard Chessboard) (res Chessboard, availableData [9][9][]int, coord Coord, ok bool) {
	// 空白格赋值标志，若中途赋值，则置为 true，需要重新获取可填值列表
	var flag bool
	var min int = 10
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			// 若当前格子存在值，则不用判断
			if chessboard[i][j] != 0 {
				continue
			}

			// 获取当前空格可填入值列表
			availableData[i][j] = GetAvailableValues(chessboard, i, j)

			// 当前空格无其他可填入值，陷入死循环，立即返回
			if len(availableData[i][j]) == 0 {
				return chessboard, availableData, coord, false
			}

			// 当前空格仅只有一个可填入值，则立即填入，将赋值标志置为 true
			if len(availableData[i][j]) == 1 {
				flag = true
				chessboard[i][j] = availableData[i][j][0]
				availableData[i][j] = nil
			}

			if len(availableData[i][j]) < min {
				min = len(availableData[i][j])
				coord = Coord{
					X: j,
					Y: i,
				}
			}
		}
	}

	// 有空白格已赋值，需重新获取可填值列表
	if flag {
		return GetChessboardAvailableValues(chessboard)
	}

	return chessboard, availableData, coord, true
}

// 获取当前坐标可填值列表。通过行、列、九宫格对 1-9 数字依次排查，获取剩余有效填入值。
func GetAvailableValues(chessboard Chessboard, i, j int) []int {
	// 已存在值，则不可再取值
	if chessboard[i][j] != 0 {
		return nil
	}

	// 需要排除的数字范围，若对应下标为 1，则需要排查，不可填入
	var exclusive [10]int
	exclusive[0] = 1

	// 九宫格对比
	for _, coord := range GridXY[GetGrid(j, i)] {
		v := chessboard[coord.Y][coord.X]
		if v != 0 {
			exclusive[v] = 1
		}
	}

	// 行对比
	for _, v := range chessboard[i] {
		if v != 0 {
			exclusive[v] = 1
		}
	}

	// 列对比
	for _, row := range chessboard {
		v := row[j]
		if v != 0 {
			exclusive[v] = 1
		}
	}

	// 获取剩余有效填入值列表
	result := make([]int, 0, 9)
	for i, v := range exclusive {
		if v == 0 {
			result = append(result, i)
		}
	}
	return result
}
