package common

import (
	"fmt"
	"testing"
	"time"
)

func TestMA(t *testing.T) {
	var tmpM []MarketData

	today := time.Now().In(time.FixedZone("CST", 28800)).Format("2006-01-02")
	yesterday := time.Now().In(time.FixedZone("CST", 28800)).Add(-24 * time.Hour).Format("2006-01-02")

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, yesterday+" 21:00:01.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, yesterday+" 21:00:11.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, yesterday+" 21:00:21.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 10, yesterday+" 21:00:31.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 100, yesterday+" 21:00:37.342 CST"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 02:04:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 02:14:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, today+" 09:00:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, today+" 09:09:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 09:11:00.342 CST"))

	result, _ := MA(tmpM, 5)

	expect := []float64{Round(float64(1.00), 2), Round(float64(3.00), 2), Round(float64(3.33333), 2), Round(float64(3.00), 2), Round(float64(3.00), 2), Round(float64(3.00), 2), Round(float64(2.60), 2), Round(float64(2.60), 2), Round(float64(3.20), 2), Round(float64(3.00), 2)}

	for i, v := range result {

		if Round(v, 2) != expect[i] {
			t.Error("TestMA() Failed. Element:", i)
			fmt.Printf("result[%d]:%v\n", i, Round(v, 2))
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestSMA(t *testing.T) {
	var tmpM []MarketData

	today := time.Now().In(time.FixedZone("CST", 28800)).Format("2006-01-02")
	yesterday := time.Now().In(time.FixedZone("CST", 28800)).Add(-24 * time.Hour).Format("2006-01-02")

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, yesterday+" 21:00:01.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, yesterday+" 21:00:11.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, yesterday+" 21:00:21.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 10, yesterday+" 21:00:31.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 100, yesterday+" 21:00:37.342 CST"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 02:04:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 02:14:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, today+" 09:00:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, today+" 09:09:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 09:11:00.342 CST"))

	result, _ := SMA(tmpM, 5)

	expect := []float64{Round(float64(3.00), 2), Round(float64(3.00), 2), Round(float64(2.60), 2), Round(float64(2.60), 2), Round(float64(3.20), 2), Round(float64(3.00), 2)}

	for i, v := range result {

		if Round(v, 2) != expect[i] {
			t.Error("TestSMA() Failed. Element:", i)
			fmt.Printf("result[%d]:%v\n", i, Round(v, 2))
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestEMA(t *testing.T) {
	var tmpM []MarketData

	yesterday := time.Now().In(time.FixedZone("CST", 28800)).Add(-24 * time.Hour).Format("2006-01-02")

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, yesterday+" 21:00:01.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, yesterday+" 21:00:11.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, yesterday+" 21:00:21.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 10, yesterday+" 21:00:31.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 100, yesterday+" 21:00:37.342 CST"))

	result, _ := EMA(tmpM, 4)

	expect := []float64{Round(float64(1.00), 2), Round(float64(2.60), 2), Round(float64(3.16), 2), Round(float64(2.696), 2), Round(float64(2.8176), 2)}

	for i, v := range result {

		if Round(v, 2) != expect[i] {
			t.Error("TestEMA() Failed. Element:", i)
			fmt.Printf("result[%d]:%v\n", i, Round(v, 2))
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestRSI(t *testing.T) {
	var tmpM []MarketData

	today := time.Now().In(time.FixedZone("CST", 28800)).Format("2006-01-02")
	yesterday := time.Now().In(time.FixedZone("CST", 28800)).Add(-24 * time.Hour).Format("2006-01-02")

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, yesterday+" 21:00:01.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, yesterday+" 21:00:11.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, yesterday+" 21:00:21.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 10, yesterday+" 21:00:31.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 100, yesterday+" 21:00:37.342 CST"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 02:04:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 02:14:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, today+" 09:00:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, today+" 09:09:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 09:11:00.342 CST"))

	result, _ := RSI(tmpM, 5)

	expect := []float64{Round(float64(-1.00), 2), Round(float64(-1.00), 2), Round(float64(-1.00), 2), Round(float64(-1.00), 2), Round(float64(-1.00), 2), Round(float64(50.00), 2), Round(float64(37.50), 2), Round(float64(50.00), 2), Round(float64(71.43), 2), Round(float64(44.44), 2)}

	for i, v := range result {

		if Round(v, 2) != expect[i] {
			t.Error("TestRSI() Failed. Element:", i)
			fmt.Printf("result[%d]:%v\n", i, Round(v, 2))
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}
}

func TestRSV(t *testing.T) {

	var tmpF []float64

	tmpF = append(tmpF, 1)
	tmpF = append(tmpF, 5)
	tmpF = append(tmpF, 4)
	tmpF = append(tmpF, 2)
	tmpF = append(tmpF, 3)

	tmpF = append(tmpF, 1)
	tmpF = append(tmpF, 3)
	tmpF = append(tmpF, 4)
	tmpF = append(tmpF, 5)
	tmpF = append(tmpF, 2)

	result, _ := RSV(tmpF, 5)

	expect := []float64{Round(float64(0), 2), Round(float64(100.00), 2), Round(float64(75.00), 2), Round(float64(25.00), 2), Round(float64(50.00), 2), Round(float64(0), 2), Round(float64(66.67), 2), Round(float64(100.00), 2), Round(float64(100.00), 2), Round(float64(25.00), 2)}

	for i, v := range result {

		if Round(v, 2) != expect[i] {
			t.Error("TestRSV() Failed. Element FLoat64:", i)
			fmt.Printf("result[%d]:%v\n", i, Round(v, 2))
			fmt.Printf("expect[%d]:%v\n", i, expect[i])
		}
	}

	var tmpM []MarketData

	today := time.Now().In(time.FixedZone("CST", 28800)).Format("2006-01-02")
	yesterday := time.Now().In(time.FixedZone("CST", 28800)).Add(-24 * time.Hour).Format("2006-01-02")

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, yesterday+" 21:00:01.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, yesterday+" 21:00:11.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, yesterday+" 21:00:21.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 10, yesterday+" 21:00:31.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 100, yesterday+" 21:00:37.342 CST"))

	tmpM = append(tmpM, GetFakeMarketData("AAA", 1, 10, today+" 02:04:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 3, 10, today+" 02:14:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 4, 10, today+" 09:00:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 5, 10, today+" 09:09:00.342 CST"))
	tmpM = append(tmpM, GetFakeMarketData("AAA", 2, 100, today+" 09:11:00.342 CST"))

	result1, _ := RSV(tmpM, 5)

	expect1 := []float64{Round(float64(0), 2), Round(float64(100.00), 2), Round(float64(75.00), 2), Round(float64(25.00), 2), Round(float64(50.00), 2), Round(float64(0), 2), Round(float64(66.67), 2), Round(float64(100.00), 2), Round(float64(100.00), 2), Round(float64(25.00), 2)}

	for i, v := range result1 {

		if Round(v, 2) != expect1[i] {
			t.Error("TestRSV() Failed. Element MarketData:", i)
			fmt.Printf("result[%d]:%v\n", i, Round(v, 2))
			fmt.Printf("expect[%d]:%v\n", i, expect1[i])
		}
	}

	var tmpK []CandleStickData
	tmpK = append(tmpK, CandleStickData{1, time.Second, "AAA", GetParsedTimeStamp(yesterday + " 21:00:00 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:01.342 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:37.342 CST"), 2, 1, 2, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{2, time.Second, "AAA", GetParsedTimeStamp(today + " 02:04:00 CST"), 1, GetParsedTimeStamp(today + " 02:04:00.342 CST"), 5, GetParsedTimeStamp(today + " 09:11:00.342 CST"), 5, 2, 3, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{3, time.Second, "AAA", GetParsedTimeStamp(today + " 09:50:00 CST"), 3, GetParsedTimeStamp(today + " 09:51:00.342 CST"), 4, GetParsedTimeStamp(today + " 10:42:00.342 CST"), 4, 3, 4, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{4, time.Second, "AAA", GetParsedTimeStamp(today + " 11:21:00 CST"), 3, GetParsedTimeStamp(today + " 11:22:00.342 CST"), 2, GetParsedTimeStamp(today + " 13:58:00.342 CST"), 3, 1, 5, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{5, time.Second, "AAA", GetParsedTimeStamp(today + " 14:37:00 CST"), 3, GetParsedTimeStamp(today + " 14:37:00.342 CST"), 3, GetParsedTimeStamp(today + " 14:59:00.342 CST"), 4, 1, 6, 100, KExpansion{}, 5})

	tmpK = append(tmpK, CandleStickData{6, time.Second, "AAA", GetParsedTimeStamp(yesterday + " 21:00:00 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:01.342 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:37.342 CST"), 3, 1, 2, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{7, time.Second, "AAA", GetParsedTimeStamp(today + " 02:04:00 CST"), 1, GetParsedTimeStamp(today + " 02:04:00.342 CST"), 3, GetParsedTimeStamp(today + " 09:11:00.342 CST"), 5, 1, 3, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{8, time.Second, "AAA", GetParsedTimeStamp(today + " 09:50:00 CST"), 3, GetParsedTimeStamp(today + " 09:51:00.342 CST"), 4, GetParsedTimeStamp(today + " 10:42:00.342 CST"), 4, 2, 4, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{9, time.Second, "AAA", GetParsedTimeStamp(today + " 11:21:00 CST"), 3, GetParsedTimeStamp(today + " 11:22:00.342 CST"), 5, GetParsedTimeStamp(today + " 13:58:00.342 CST"), 5, 1, 5, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{10, time.Second, "AAA", GetParsedTimeStamp(today + " 14:37:00 CST"), 3, GetParsedTimeStamp(today + " 14:37:00.342 CST"), 2, GetParsedTimeStamp(today + " 14:59:00.342 CST"), 2, 1, 6, 100, KExpansion{}, 5})

	result2, _ := RSV(tmpK, 5)

	expect2 := []float64{Round(float64(0), 2), Round(float64(100.00), 2), Round(float64(75.00), 2), Round(float64(25.00), 2), Round(float64(50.00), 2), Round(float64(0), 2), Round(float64(50.00), 2), Round(float64(75.00), 2), Round(float64(100.00), 2), Round(float64(25.00), 2)}

	for i, v := range result2 {

		if Round(v, 2) != expect2[i] {
			t.Error("TestRSV() Failed. Element CandleStickData:", i)
			fmt.Printf("result[%d]:%v\n", i, Round(v, 2))
			fmt.Printf("expect[%d]:%v\n", i, expect2[i])
		}
	}
}

func TestKD(t *testing.T) {

	today := time.Now().In(time.FixedZone("CST", 28800)).Format("2006-01-02")
	yesterday := time.Now().In(time.FixedZone("CST", 28800)).Add(-24 * time.Hour).Format("2006-01-02")

	var tmpK []CandleStickData
	tmpK = append(tmpK, CandleStickData{1, time.Second, "AAA", GetParsedTimeStamp(yesterday + " 21:00:00 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:01.342 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:37.342 CST"), 2, 1, 2, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{2, time.Second, "AAA", GetParsedTimeStamp(today + " 02:04:00 CST"), 1, GetParsedTimeStamp(today + " 02:04:00.342 CST"), 5, GetParsedTimeStamp(today + " 09:11:00.342 CST"), 5, 2, 3, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{3, time.Second, "AAA", GetParsedTimeStamp(today + " 09:50:00 CST"), 3, GetParsedTimeStamp(today + " 09:51:00.342 CST"), 4, GetParsedTimeStamp(today + " 10:42:00.342 CST"), 4, 3, 4, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{4, time.Second, "AAA", GetParsedTimeStamp(today + " 11:21:00 CST"), 3, GetParsedTimeStamp(today + " 11:22:00.342 CST"), 2, GetParsedTimeStamp(today + " 13:58:00.342 CST"), 3, 1, 5, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{5, time.Second, "AAA", GetParsedTimeStamp(today + " 14:37:00 CST"), 3, GetParsedTimeStamp(today + " 14:37:00.342 CST"), 3, GetParsedTimeStamp(today + " 14:59:00.342 CST"), 4, 1, 6, 100, KExpansion{}, 5})

	tmpK = append(tmpK, CandleStickData{6, time.Second, "AAA", GetParsedTimeStamp(yesterday + " 21:00:00 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:01.342 CST"), 1, GetParsedTimeStamp(yesterday + " 21:00:37.342 CST"), 3, 1, 2, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{7, time.Second, "AAA", GetParsedTimeStamp(today + " 02:04:00 CST"), 1, GetParsedTimeStamp(today + " 02:04:00.342 CST"), 3, GetParsedTimeStamp(today + " 09:11:00.342 CST"), 5, 1, 3, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{8, time.Second, "AAA", GetParsedTimeStamp(today + " 09:50:00 CST"), 3, GetParsedTimeStamp(today + " 09:51:00.342 CST"), 4, GetParsedTimeStamp(today + " 10:42:00.342 CST"), 4, 2, 4, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{9, time.Second, "AAA", GetParsedTimeStamp(today + " 11:21:00 CST"), 3, GetParsedTimeStamp(today + " 11:22:00.342 CST"), 5, GetParsedTimeStamp(today + " 13:58:00.342 CST"), 5, 1, 5, 100, KExpansion{}, 5})
	tmpK = append(tmpK, CandleStickData{10, time.Second, "AAA", GetParsedTimeStamp(today + " 14:37:00 CST"), 3, GetParsedTimeStamp(today + " 14:37:00.342 CST"), 2, GetParsedTimeStamp(today + " 14:59:00.342 CST"), 2, 1, 6, 100, KExpansion{}, 5})

	result, _ := KD(tmpK, 5)

	expect := []KDData{KDData{Round(float64(0), 2), Round(float64(0), 2)},
		KDData{Round(float64(33.33), 2), Round(float64(11.11), 2)},
		KDData{Round(float64(47.22), 2), Round(float64(23.15), 2)},
		KDData{Round(float64(39.81), 2), Round(float64(28.70), 2)},
		KDData{Round(float64(43.21), 2), Round(float64(33.54), 2)},
		KDData{Round(float64(28.81), 2), Round(float64(31.96), 2)},
		KDData{Round(float64(35.87), 2), Round(float64(33.26), 2)},
		KDData{Round(float64(48.91), 2), Round(float64(38.48), 2)},
		KDData{Round(float64(65.94), 2), Round(float64(47.64), 2)},
		KDData{Round(float64(52.30), 2), Round(float64(49.19), 2)}}

	for i, v := range result {

		if Round(v.K, 2) != expect[i].K || Round(v.D, 2) != expect[i].D {
			t.Error("TestRSV() Failed. Element CandleStickData:", i)
			fmt.Printf("result[%d]:[K]%v,[D]%v\n", i, Round(v.K, 2), Round(v.D, 2))
			fmt.Printf("expect[%d]:[K]%v,[D]%v\n", i, expect[i].K, expect[i].D)
		}
	}
}
