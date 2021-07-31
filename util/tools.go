package util

import (
	"math"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func GetCurrentFilename() string {
	_, current_file, _, _ := runtime.Caller(1)
	baseFile := filepath.Base(current_file)
	return strings.TrimSuffix(baseFile, filepath.Ext(baseFile))
}

func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

func AppendString(arr []string, str *string) []string {
	if str != nil && *str != "" {
		arr = append(arr, *str)
	}
	return arr
}

// adapted from https://www.socketloop.com/tutorials/golang-byte-format-example
func ByteFormat(inputNum int64, precision int) string {

	if precision <= 0 {
		precision = 1
	}

	var unit string
	var returnVal float64

	floatInputNum := float64(inputNum)
	if inputNum >= 1000000000000000 {
		returnVal = RoundUp((floatInputNum / 1125899906842624), precision)
		unit = "PB" // petabyte
	} else if floatInputNum >= 1000000000000 {
		returnVal = RoundUp((floatInputNum / 1099511627776), precision)
		unit = "TB" // terrabyte
	} else if floatInputNum >= 1000000000 {
		returnVal = RoundUp((floatInputNum / 1073741824), precision)
		unit = "GB" // gigabyte
	} else if floatInputNum >= 1000000 {
		returnVal = RoundUp((floatInputNum / 1048576), precision)
		unit = "MB" // megabyte
	} else if floatInputNum >= 1000 {
		returnVal = RoundUp((floatInputNum / 1024), precision)
		unit = "KB" // kilobyte
	} else {
		returnVal = floatInputNum
		unit = "B" // byte
	}

	return strconv.FormatFloat(returnVal, 'f', -1, 64) + " " + unit

}

func ReverseString(s string) string {
	var sb strings.Builder
	runes := []rune(s)
	for i := len(runes) - 1; 0 <= i; i-- {
		sb.WriteRune(runes[i])
	}
	return sb.String()
}

func ReplaceRight(s string, old string, new string, n int) string {
	return ReverseString(strings.Replace(ReverseString(s), ReverseString(old), ReverseString(new), 1))
}
