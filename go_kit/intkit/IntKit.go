package intkit

import "math"

// @Title 判断多个整形是否为0
// @Description
// @param strs
// usage:
//	IntIsZero(0) return false
func IntIsZero(ints ...int) bool {
	if len(ints) == 0 {
		return false
	}
	for _, v := range ints {
		if v != 0 {
			return false
		}
	}
	return true
}

// 小数转整数四舍五入
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// 小数保留指定小数位数
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
