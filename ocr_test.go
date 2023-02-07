package sudoku

import (
	"testing"
)

var _testChessboard = Chessboard{
	{0, 5, 9, 3, 0, 2, 0, 0, 0},
	{0, 2, 0, 4, 1, 0, 8, 0, 0},
	{0, 0, 0, 7, 0, 0, 6, 0, 0},
	{1, 3, 0, 0, 6, 7, 9, 0, 4},
	{4, 0, 0, 0, 5, 0, 0, 0, 7},
	{9, 6, 7, 0, 4, 8, 5, 3, 2},
	{2, 7, 8, 6, 3, 0, 4, 5, 9},
	{3, 9, 6, 5, 0, 0, 0, 1, 8},
	{0, 1, 0, 8, 2, 0, 0, 0, 0},
}

func TestReadSudoImg(t *testing.T) {
	res, err := ReadSudoImg("sudu.png")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if res != _testChessboard {
		t.Errorf("期望九宫格：\n%v\n实际九宫格：\n%v\n", _testChessboard, res)
		return
	}

	r, ok := Sudoku(res)
	if !ok {
		t.Errorf("获取九宫格正确结果失败\n")
		return
	}
	t.Logf("源九宫格：\n%v\n补全后的九宫格：\n%v\n", res, r)
}
