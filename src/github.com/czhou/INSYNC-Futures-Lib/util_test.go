package common

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

/*
func TestGetTradingTimeSlots(t *testing.T) {
	//测试品种SH_C1，有夜盘凌晨2点30收盘
	//	SH_BID_DAY_START      = "08:55:00"
	//	SH_BID_DAY_END        = "09:00:00"
	//	SH_TRADE_DAY_START_P1 = "09:00:00"
	//	SH_TRADE_DAY_END_P1   = "10:15:00"
	//	SH_TRADE_DAY_START_P2 = "10:30:00"
	//	SH_TRADE_DAY_END_P2   = "11:30:00"
	//	SH_TRADE_DAY_START_P3 = "13:30:00"
	//	SH_TRADE_DAY_END_P3   = "15:00:00"

	//	SH_C1_BID_NIGHT_START      = "20:55:00"
	//	SH_C1_BID_NIGHT_END        = "21:00:00"
	//	SH_C1_TRADE_NIGHT_START_P1 = "21:00:00"
	//	SH_C1_TRADE_NIGHT_END_P1   = "02:30:00" //次日

	result, err := GetTradingTimeSlots(SH_C1)
	assert.Equal(t, nil, err)

	today := time.Now().In(time.FixedZone("CST", 28800)).Format("2006-01-02")
	yesterday := time.Now().In(time.FixedZone("CST", 28800)).Add(-24 * time.Hour).Format("2006-01-02")

	t1, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, yesterday+" 21:00:00 CST")
	t2, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 02:30:00 CST")
	t3, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 09:00:00 CST")
	t4, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 10:15:00 CST")
	t5, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 10:30:00 CST")
	t6, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 11:30:00 CST")
	t7, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 13:30:00 CST")
	t8, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 15:00:00 CST")

	expect := []time.Time{t1, t2, t3, t4, t5, t6, t7, t8}

	for i, v := range result {
		if !v.Equal(expect[i]) {
			t.Error("TestGetTradingTimeSlots() Failed. Initial timeslots is not as expected, Element:", i)
			fmt.Printf("result[%d]:%v\n", i, result[i])
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}
*/

func TestGetTradingTimeSlots(t *testing.T) {

	result, err := GetTradingTimeSlots("ni")
	assert.Equal(t, nil, err)

	timeNow := time.Now().In(time.FixedZone("CST", 28800))
	//如果是下午7点以后调用该函数，则返回第二天的时间
	if timeNow.Hour() > 19 {
		timeNow = timeNow.Add(time.Hour * 24)
	}

	today := timeNow.Format("2006-01-02")
	yesterday := timeNow.Add(-24 * time.Hour).Format("2006-01-02")

	t1, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, yesterday+" 21:00:00 CST")
	t2, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 1:00:00 CST")
	t3, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 09:00:00 CST")
	t4, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 10:15:00 CST")
	t5, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 10:30:00 CST")
	t6, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 11:30:00 CST")
	t7, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 13:30:00 CST")
	t8, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 15:00:00 CST")

	expect := []time.Time{t1, t2, t3, t4, t5, t6, t7, t8}
	//expect := []time.Time{t3, t4, t5, t6, t7, t8}
	//expect := []time.Time{t3, t4, t5, t6}

	for i, v := range result {
		if !v.Equal(expect[i]) {
			t.Error("TestGetTradingTimeSlots() Failed. Initial timeslots is not as expected, Element:", i)
			fmt.Printf("result[%d]:%v\n", i, result[i])
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestGetIntervalTimeSlots(t *testing.T) {
	//测试品种有夜盘凌晨2点30收盘
	//K线周期15分钟和38分钟

	timeNow := time.Now().In(time.FixedZone("CST", 28800))
	//如果是下午7点以后调用该函数，则返回第二天的时间
	if timeNow.Hour() > 19 {
		timeNow = timeNow.Add(time.Hour * 24)
	}

	today := timeNow.Format("2006-01-02")
	yesterday := timeNow.Add(-24 * time.Hour).Format("2006-01-02")

	result, _ := GetIntervalTimeSlots("ag", time.Minute*38)

	t1, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, yesterday+" 21:00:00 CST")
	t2, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, yesterday+" 21:38:00 CST")
	t3, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, yesterday+" 22:16:00 CST")
	t4, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, yesterday+" 22:54:00 CST")
	t5, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, yesterday+" 23:32:00 CST")
	t6, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 00:10:00 CST")
	t7, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 00:48:00 CST")
	t8, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 01:26:00 CST")
	t9, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 02:04:00 CST")
	t10, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 09:12:00 CST")
	t11, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 09:50:00 CST")
	t12, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 10:43:00 CST")
	t13, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 11:21:00 CST")
	t14, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 13:59:00 CST")
	t15, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 14:37:00 CST")
	t16, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, today+" 15:00:00 CST")

	expect := []time.Time{t1, t2, t3, t4, t5, t6, t7, t8, t9, t10, t11, t12, t13, t14, t15, t16}

	for i, v := range result {

		if !v.Equal(expect[i]) {
			t.Error("TestGetIntervalTimeSlots() Failed. Candlestick timeslots is not as expected, Element:", i)
			fmt.Printf("result[%d]:%v\n", i, result[i])
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestGetCandleStickData(t *testing.T) {
	var tmpM []MarketData

	timeNow := time.Now().In(time.FixedZone("CST", 28800))
	//如果是下午7点以后调用该函数，则返回第二天的时间
	if timeNow.Hour() > 19 {
		timeNow = timeNow.Add(time.Hour * 24)
	}

	today := timeNow.Format("2006-01-02")
	yesterday := timeNow.Add(-24 * time.Hour).Format("2006-01-02")

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, yesterday+" 21:00:01.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, yesterday+" 21:00:11.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, yesterday+" 21:00:21.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 10, yesterday+" 21:00:31.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 100, yesterday+" 21:00:37.342"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 02:04:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 02:14:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, today+" 09:00:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, today+" 09:09:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 09:11:00.342"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 09:51:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, today+" 09:54:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 10:00:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, today+" 10:31:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 10:42:00.342"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 11:22:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, today+" 11:29:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 13:00:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, today+" 13:21:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 13:58:00.342"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 14:37:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, today+" 14:38:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 14:50:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, today+" 14:55:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 14:59:00.342"))

	kTimeSlots, _ := GetIntervalTimeSlots("ag", time.Minute*38)

	var result []CandleStickData

	for _, v := range tmpM {
		result = GetCandleStickData(kTimeSlots, result, v)
	}

	var expect []CandleStickData
	expect = append(expect, CandleStickData{1, time.Second, "AAA", GetParsedTimeStamp(yesterday + " 21:00:00 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:01.342 CST"), 3, GetParsedTimeStamp(yesterday + " 21:00:37.342 CST"), 5, 1, 2, 100, KExpansion{}, 5})
	expect = append(expect, CandleStickData{2, time.Second, "AAA", GetParsedTimeStamp(today + " 02:04:00 CST"), 1, GetParsedTimeStamp(today + " 02:04:00.342 CST"), 2, GetParsedTimeStamp(today + " 09:11:00.342 CST"), 5, 1, 3, 100, KExpansion{}, 5})
	expect = append(expect, CandleStickData{3, time.Second, "AAA", GetParsedTimeStamp(today + " 09:50:00 CST"), 3, GetParsedTimeStamp(today + " 09:51:00.342 CST"), 2, GetParsedTimeStamp(today + " 10:42:00.342 CST"), 5, 1, 4, 100, KExpansion{}, 5})
	expect = append(expect, CandleStickData{4, time.Second, "AAA", GetParsedTimeStamp(today + " 11:21:00 CST"), 3, GetParsedTimeStamp(today + " 11:22:00.342 CST"), 2, GetParsedTimeStamp(today + " 13:58:00.342 CST"), 5, 1, 5, 100, KExpansion{}, 5})
	expect = append(expect, CandleStickData{5, time.Second, "AAA", GetParsedTimeStamp(today + " 14:37:00 CST"), 3, GetParsedTimeStamp(today + " 14:37:00.342 CST"), 2, GetParsedTimeStamp(today + " 14:59:00.342 CST"), 5, 1, 6, 100, KExpansion{}, 5})

	for i, _ := range result {
		if !IsEqual(result[i], expect[i]) {
			t.Error("TestGetCandleStickData() Failed. Candlestick data is not as expected, Element:", i)
			fmt.Printf("result[%d]:%v\n", i, result[i])
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestGetComplexCandleStickData(t *testing.T) {
	var tmpM []MarketData

	timeNow := time.Now().In(time.FixedZone("CST", 28800))
	//如果是下午7点以后调用该函数，则返回第二天的时间
	if timeNow.Hour() > 19 {
		timeNow = timeNow.Add(time.Hour * 24)
	}

	today := timeNow.Format("2006-01-02")
	yesterday := timeNow.Add(-24 * time.Hour).Format("2006-01-02")

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, yesterday+" 21:00:01.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 20, yesterday+" 21:00:11.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 30, yesterday+" 21:00:21.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 40, yesterday+" 21:00:31.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 100, yesterday+" 21:00:37.342"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 02:04:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 20, today+" 02:14:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 30, today+" 09:00:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 40, today+" 09:09:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 09:11:00.342"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 09:51:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 20, today+" 09:54:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 30, today+" 10:00:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 40, today+" 10:31:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 10:42:00.342"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 11:22:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 20, today+" 11:29:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 30, today+" 13:00:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 40, today+" 13:21:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 13:58:00.342"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 14:37:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 20, today+" 14:38:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 30, today+" 14:50:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 40, today+" 14:55:00.342"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 14:59:00.342"))

	kTimeSlots, _ := GetIntervalTimeSlots("ag", time.Minute*38)

	var result []CandleStickData

	for _, v := range tmpM {
		result = GetComplexCandleStickData(kTimeSlots, result, v)
	}

	var expect []CandleStickData
	expect = append(expect, CandleStickData{1, time.Minute * 38, "AAA", GetParsedTimeStamp(yesterday + " 21:00:00 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:01.342 CST"), 3, GetParsedTimeStamp(yesterday + " 21:00:37.342 CST"), 5, 1, 2, 100, KExpansion{RoundToIntDigit(610), 200, Round(3.05, 2)}, 5})
	expect = append(expect, CandleStickData{2, time.Minute * 38, "AAA", GetParsedTimeStamp(today + " 02:04:00 CST"), 1, GetParsedTimeStamp(today + " 02:04:00.342 CST"), 2, GetParsedTimeStamp(today + " 09:11:00.342 CST"), 5, 1, 3, 100, KExpansion{RoundToIntDigit(590), 200, Round(2.95, 2)}, 5})
	expect = append(expect, CandleStickData{3, time.Minute * 38, "AAA", GetParsedTimeStamp(today + " 09:50:00 CST"), 3, GetParsedTimeStamp(today + " 09:51:00.342 CST"), 2, GetParsedTimeStamp(today + " 10:42:00.342 CST"), 5, 1, 4, 100, KExpansion{RoundToIntDigit(540), 200, Round(2.70, 2)}, 5})
	expect = append(expect, CandleStickData{4, time.Minute * 38, "AAA", GetParsedTimeStamp(today + " 11:21:00 CST"), 3, GetParsedTimeStamp(today + " 11:22:00.342 CST"), 2, GetParsedTimeStamp(today + " 13:58:00.342 CST"), 5, 1, 5, 100, KExpansion{RoundToIntDigit(540), 200, Round(2.70, 2)}, 5})
	expect = append(expect, CandleStickData{5, time.Minute * 38, "AAA", GetParsedTimeStamp(today + " 14:37:00 CST"), 3, GetParsedTimeStamp(today + " 14:37:00.342 CST"), 2, GetParsedTimeStamp(today + " 14:59:00.342 CST"), 5, 1, 6, 100, KExpansion{RoundToIntDigit(540), 200, Round(2.70, 2)}, 5})

	for i, _ := range result {
		if !IsEqual(result[i], expect[i]) {
			t.Error("TestGetCandleStickData() Failed. Candlestick data is not as expected, Element:", i)
			fmt.Printf("result[%d]:%v\n", i, result[i])
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func GetFakeMarketData(instID string, lastPrice float64, volume int64, timeS string) MarketData {
	var tmp MarketData
	tmp.InstrumentID = instID
	tmp.LastPrice = lastPrice
	tmp.Volume = volume
	tmp.TimeStamp = timeS

	return tmp
}

func GetParsedTimeStamp(t string) time.Time {
	tmp, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, t)
	return tmp
}

func IsEqual(p1, p2 interface{}) bool {

	switch p1.(type) {

	case CandleStickData:

		v1 := p1.(CandleStickData)
		v2 := p2.(CandleStickData)

		if v1.KNum == v2.KNum && v1.InstrumentID == v2.InstrumentID && v1.KTime.Equal(v2.KTime) && v1.OpenPrice == v2.OpenPrice && v1.OpenTimestamp.Equal(v2.OpenTimestamp) && v1.ClosePrice == v2.ClosePrice && v1.CloseTimestamp.Equal(v2.CloseTimestamp) && v1.HighestPrice == v2.HighestPrice && v1.LowestPrice == v2.LowestPrice && v1.Count == v2.Count && v1.Expansion.PriceVolumeTotal == v2.Expansion.PriceVolumeTotal && v1.Expansion.VolumeTotal == v2.Expansion.VolumeTotal && v1.Expansion.Vwap == v2.Expansion.Vwap {
			return true
		} else {
			return false
		}

	default:
		return false
	}

}
