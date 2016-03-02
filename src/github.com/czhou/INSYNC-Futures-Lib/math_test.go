package common

import (
	"fmt"
	"testing"
)

func TestRoundToIntDigit(t *testing.T) {
	floatArray := []float64{-1.5, -1.2, 0, 1.2, 1.5}
	expect := []float64{-2, -1, 0, 1, 2}

	for i, v := range floatArray {
		if RoundToIntDigit(v) != expect[i] {
			t.Error("TestEMA() Failed. Element:", i)
			fmt.Printf("result[%d]:%v\n", i, v)
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestRound(t *testing.T) {
	floatArray := []float64{-1.524, -1.526, 0, 1.234, 1.236}
	expect := []float64{-1.52, -1.53, 0, 1.23, 1.24}

	for i, v := range floatArray {
		if Round(v, 2) != expect[i] {
			t.Error("TestRound() Failed. Element:", i)
			fmt.Printf("result[%d]:%v\n", i, Round(v, 2))
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestAvg(t *testing.T) {
	floatArray := []float64{-1, 0, 1, 2, 3, 4, 5, 6}
	expect := float64(2.5)

	result, _ := Avg(floatArray)

	if Round(result, 1) != expect {
		t.Error("TestAvg() Failed.")
		fmt.Printf("result:%v\n", Round(result, 1))
		fmt.Printf("expect:%v\n", expect)
	}
}

func TestMax(t *testing.T) {
	floatArray := []float64{-1, 0, 1, 2, 3, 4, 5, 6}
	expect := float64(6)

	i, result, _ := Max(floatArray)

	if i != 7 || result != expect {
		t.Error("TestMax() Failed.")
		fmt.Printf("result[%d]:%v\n", i, result)
		fmt.Printf("expect[7]:%v\n", expect)
	}
}

func TestMin(t *testing.T) {
	floatArray := []float64{-1, 0, 1, 2, 3, 4, 5, 6}
	expect := float64(-1)

	i, result, _ := Min(floatArray)

	if i != 0 || result != expect {
		t.Error("TestMin() Failed.")
		fmt.Printf("result[%d]:%v\n", i, result)
		fmt.Printf("expect[0]:%v\n", expect)
	}
}
