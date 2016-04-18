package common

import (
	"math"
)

//保留整数，四舍五入
func RoundToIntDigit(value float64) float64 {
	if value < 0 {
		return math.Ceil(value - 0.5)
	}
	return math.Floor(value + 0.5)
}

//保留小数点后多少位，四舍五入
func Round(value float64, digits int) float64 {
	shift := math.Pow(10, float64(digits))
	return RoundToIntDigit(value*shift) / shift
}

//求算数平均数
//支持传入所有int和float类型的数组
//返回参数：平均值，错误
func Avg(array interface{}) (avgVal float64, err error) {

	var sum float64 = 0

	switch val := array.(type) {

	case []int:

		if len(val) == 0 {
			return 0, Error{"Input array for function Avg is empty."}
		}

		count := 0
		for i, v := range val {
			sum = sum + float64(v)
			count = i + 1
		}
		return sum / float64(count), nil

	case []int8:

		if len(val) == 0 {
			return 0, Error{"Input array for function Avg is empty."}
		}

		count := 0
		for i, v := range val {
			sum = sum + float64(v)
			count = i + 1
		}
		return sum / float64(count), nil

	case []int16:

		if len(val) == 0 {
			return 0, Error{"Input array for function Avg is empty."}
		}

		count := 0
		for i, v := range val {
			sum = sum + float64(v)
			count = i + 1
		}
		return sum / float64(count), nil

	case []int32:

		if len(val) == 0 {
			return 0, Error{"Input array for function Avg is empty."}
		}

		count := 0
		for i, v := range val {
			sum = sum + float64(v)
			count = i + 1
		}
		return sum / float64(count), nil

	case []int64:

		if len(val) == 0 {
			return 0, Error{"Input array for function Avg is empty."}
		}

		count := 0
		for i, v := range val {
			sum = sum + float64(v)
			count = i + 1
		}
		return sum / float64(count), nil

	case []float32:

		if len(val) == 0 {
			return 0, Error{"Input array for function Avg is empty."}
		}

		count := 0
		for i, v := range val {
			sum = sum + float64(v)
			count = i + 1
		}
		return sum / float64(count), nil

	case []float64:

		if len(val) == 0 {
			return 0, Error{"Input array for function Avg is empty."}
		}

		count := 0
		for i, v := range val {
			sum = sum + v
			count = i + 1
		}
		return sum / float64(count), nil

	default:
		return 0, Error{"Input type is not supported by function Avg."}
	}
}

//求最大值
//支持传入所有int和float类型的数组
//返回参数：元素位置，最大值，错误
func Max(array interface{}) (index int, maxVal interface{}, err error) {

	switch val := array.(type) {

	case []int:

		var max int
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		max = val[0]

		for i := 1; i < len(val); i++ {
			if max < val[i] {
				max = val[i]
				index = i
			}
		}

		return index, max, nil

	case []int8:

		var max int8
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		max = val[0]

		for i := 1; i < len(val); i++ {
			if max < val[i] {
				max = val[i]
				index = i
			}
		}

		return index, max, nil

	case []int16:

		var max int16
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		max = val[0]

		for i := 1; i < len(val); i++ {
			if max < val[i] {
				max = val[i]
				index = i
			}
		}

		return index, max, nil

	case []int32:

		var max int32
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		max = val[0]

		for i := 1; i < len(val); i++ {
			if max < val[i] {
				max = val[i]
				index = i
			}
		}

		return index, max, nil

	case []int64:

		var max int64
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		max = val[0]

		for i := 1; i < len(val); i++ {
			if max < val[i] {
				max = val[i]
				index = i
			}
		}

		return index, max, nil

	case []float32:

		var max float32
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		max = val[0]

		for i := 1; i < len(val); i++ {
			if max < val[i] {
				max = val[i]
				index = i
			}
		}

		return index, max, nil

	case []float64:

		var max float64
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		max = val[0]

		for i := 1; i < len(val); i++ {
			if max < val[i] {
				max = val[i]
				index = i
			}
		}

		return index, max, nil

	default:
		return -1, 0, Error{"Input type is not supported by function Avg."}
	}
}

//求最小值
//支持传入所有int和float类型的数组
//返回参数：元素位置，最小值，错误
func Min(array interface{}) (index int, minVal interface{}, err error) {

	switch val := array.(type) {

	case []int:

		var min int
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		min = val[0]

		for i := 1; i < len(val); i++ {
			if min > val[i] {
				min = val[i]
				index = i
			}
		}

		return index, min, nil

	case []int8:

		var min int8
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		min = val[0]

		for i := 1; i < len(val); i++ {
			if min > val[i] {
				min = val[i]
				index = i
			}
		}

		return index, min, nil

	case []int16:

		var min int16
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		min = val[0]

		for i := 1; i < len(val); i++ {
			if min > val[i] {
				min = val[i]
				index = i
			}
		}

		return index, min, nil

	case []int32:

		var min int32
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		min = val[0]

		for i := 1; i < len(val); i++ {
			if min > val[i] {
				min = val[i]
				index = i
			}
		}

		return index, min, nil

	case []int64:

		var min int64
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		min = val[0]

		for i := 1; i < len(val); i++ {
			if min > val[i] {
				min = val[i]
				index = i
			}
		}

		return index, min, nil

	case []float32:

		var min float32
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		min = val[0]

		for i := 1; i < len(val); i++ {
			if min > val[i] {
				min = val[i]
				index = i
			}
		}

		return index, min, nil

	case []float64:

		var min float64
		var index int

		if len(val) == 0 {
			return -1, 0, Error{"Input array for function Avg is empty."}
		}
		if len(val) == 1 {
			return 0, val[0], nil
		}

		min = val[0]

		for i := 1; i < len(val); i++ {
			if min > val[i] {
				min = val[i]
				index = i
			}
		}

		return index, min, nil

	default:
		return -1, 0, Error{"Input type is not supported by function Avg."}
	}
}
