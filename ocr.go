package sudoku

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type (
	TextRegion [4][2]float64 // 文本坐标，分别对应左上、右上、右下、左下的坐标点
	TextObj    struct {
		Confidence float64    `json:"confidence,omitempty"` // 文本
		Text       string     `json:"text,omitempty"`       // 文本识别内容
		TextRegion TextRegion `json:"text_region,omitempty"`
	}

	Bbox [4]float64 // 检测出的表格左上角和右下角坐标
)

// paddleHub 部署：https://hub.docker.com/r/scholarli/paddlehub
const paddleHubDomain = "http://127.0.0.1:9000"

func ReadSudoImg(file string) (chessboard Chessboard, err error) {
	b, err := Base64ImageFile(file)
	if err != nil {
		return chessboard, err
	}

	imgBase64 := string(b)
	objs, err := paddleOCR(imgBase64)
	if err != nil {
		return chessboard, err
	}

	bbox, err := paddleLayout(imgBase64)
	if err != nil {
		return chessboard, err
	}

	avgX := (bbox[2] - bbox[0]) / 9
	avgY := (bbox[3] - bbox[1]) / 9
	for _, obj := range objs {
		num, err := strconv.ParseUint(strings.TrimSpace(obj.Text), 10, 64)
		if err != nil {
			continue
		}
		coord := Coord{
			X: int((obj.TextRegion[0][0] - bbox[0]) / avgX),
			Y: int((obj.TextRegion[1][1] - bbox[1]) / avgY),
		}
		chessboard[coord.Y][coord.X] = int(num)
	}
	return chessboard, nil
}

func paddleOCR(imgBase64 string) ([]TextObj, error) {
	body, err := paddleAPI(paddleHubDomain+"/predict/ocr_system", imgBase64)
	if err != nil {
		return nil, err
	}

	var res struct {
		Results [][]TextObj `json:"results"`
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Results) > 0 {
		return res.Results[0], nil
	}
	return nil, fmt.Errorf("not found any txt")
}

func paddleLayout(imgBase64 string) (Bbox, error) {
	body, err := paddleAPI(paddleHubDomain+"/predict/structure_layout", imgBase64)
	if err != nil {
		return Bbox{}, err
	}

	var res struct {
		Results []struct {
			Layout []struct {
				Bbox Bbox `json:"bbox"`
			} `json:"layout"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return Bbox{}, err
	}
	if len(res.Results) > 0 && len(res.Results[0].Layout) > 0 {
		return res.Results[0].Layout[0].Bbox, nil
	}
	return Bbox{}, fmt.Errorf("not found bbox")
}

func paddleAPI(url string, imgBase64 string) ([]byte, error) {
	body := map[string]interface{}{
		"images": []string{imgBase64},
	}

	header := map[string]string{
		"Content-Type": "application/json",
	}

	code, res, err := Post(url, body, header)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("code != 200")
	}
	return res, nil
}
