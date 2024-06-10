package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

func GetImgInfoFromByte(base64Str string) (int, int, int, error) {
	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		logx.Errorf("decode base64 error1:%v", err)
		return 0, 0, 0, err
	}
	img, _, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		logx.Errorf("decode base64 error2:%v", err)
		return 0, 0, 0, err
	}
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	fileSize := len(decoded)
	return width, height, fileSize, nil
}

// InterfaceToInt 将interface{} 转换成 int
// 强制转换，忽略精度丢失; 忽略忽略错误，错误返回0
func InterfaceToInt(key interface{}) int {
	if key == nil { // nil返回零值，根据需求添加
		return 0
	}
	var ret int
	switch key := key.(type) {
	case string:
		tmp, _ := strconv.ParseFloat(key, 32)
		ret = int(tmp)
	case int:
		ret = int(key)
	case int8:
		ret = int(key)
	case int16:
		ret = int(key)
	case int32:
		ret = int(key)
	case int64:
		ret = int(key)
	case uint:
		ret = int(key)
	case uint8:
		ret = int(key)
	case uint16:
		ret = int(key)
	case uint32:
		ret = int(key)
	case uint64:
		ret = int(key)
	case float32:
		ret = int(key)
	case float64:
		ret = int(key)
	default:
	}
	return ret
}
