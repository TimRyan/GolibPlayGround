package common

import (
	"fmt"
	"time"
)

//字符串截取
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//一次返回所有K线时间段
//直接传入品种代码instCode，自动识别交易时间
func GetIntervalTimeSlots(instCode string, duration time.Duration, timeNow time.Time) ([]time.Time, error) {

	var category int //品种的交易时间类型

	switch instCode {
	//大商所品种
	case "a", "b", "bb", "c", "fb", "i", "j", "jd", "jm", "m", "p", "v", "y":
		category = NIGHT_C2
	//大商所品种
	case "cs", "l", "pp":
		category = DAY_C1
	//上期所品种
	case "ag", "au":
		category = NIGHT_C4
	//上期所品种
	case "al", "bu", "cu", "hc", "ni", "pb", "rb", "sn", "zn":
		category = NIGHT_C3
	//上期所品种
	case "ru":
		category = NIGHT_C1
	//上期所品种
	case "fu", "wr":
		category = DAY_C1
	//郑商所品种
	case "CF", "FG", "JR", "LR", "MA", "OI", "PM", "RI", "RM", "RS", "SF", "SM", "SR", "TA", "TC", "WH", "ZC":
		category = NIGHT_C2
	//中金所品种
	case "IC", "IH", "IF", "T", "TF":
		category = ZJ_DAY

	default:
		category = DAY_C1
	}

	var timeslots []time.Time
	var nonTradeDuration []time.Duration

	tradingTimeSlots, err := GetTradingTimeSlots(instCode, timeNow)

	if err != nil {
		return []time.Time{time.Now()}, err
	}

	if duration <= 0 {
		return []time.Time{tradingTimeSlots[0], tradingTimeSlots[len(tradingTimeSlots)-1]}, nil
	}

	switch category {

	case NIGHT_C4, NIGHT_C3, NIGHT_C2, NIGHT_C1:

		//交易时间间隔
		nonTradeDuration = append(nonTradeDuration, tradingTimeSlots[2].Sub(tradingTimeSlots[1]))
		nonTradeDuration = append(nonTradeDuration, tradingTimeSlots[4].Sub(tradingTimeSlots[3]))
		nonTradeDuration = append(nonTradeDuration, tradingTimeSlots[6].Sub(tradingTimeSlots[5]))

		//添加起始时间
		timeslots = append(timeslots, tradingTimeSlots[0])

		//计算每个周期时间段
		for {

			tTmp := timeslots[len(timeslots)-1].Add(duration)

			if !tTmp.After(tradingTimeSlots[1]) {
				if tTmp.Equal(tradingTimeSlots[1]) {
					timeslots = append(timeslots, tradingTimeSlots[2])
				} else {
					timeslots = append(timeslots, tTmp)
				}
				continue
			} else if !tTmp.After(tradingTimeSlots[3]) {

				if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[2]) {
					timeslots = append(timeslots, tTmp)
					continue
				} else {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[0]))
					continue
				}

			} else if !tTmp.After(tradingTimeSlots[5]) {

				if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[4]) {
					timeslots = append(timeslots, tTmp)
					continue
				} else if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[2]) {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[1]))
					continue
				} else {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[1]).Add(nonTradeDuration[0]))
					continue
				}

			} else if !tTmp.After(tradingTimeSlots[7]) {
				if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[6]) {
					timeslots = append(timeslots, tTmp)
					continue
				} else if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[4]) {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[2]))
					continue
				} else if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[2]) {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[2]).Add(nonTradeDuration[1]))
					continue
				} else {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[2]).Add(nonTradeDuration[1]).Add(nonTradeDuration[0]))
					continue
				}

			} else {
				if !timeslots[len(timeslots)-1].Equal(tradingTimeSlots[7]) {
					timeslots = append(timeslots, tradingTimeSlots[7])
				}
				break
			}

		}

	case DAY_C1:

		//交易时间间隔
		nonTradeDuration = append(nonTradeDuration, tradingTimeSlots[2].Sub(tradingTimeSlots[1]))
		nonTradeDuration = append(nonTradeDuration, tradingTimeSlots[4].Sub(tradingTimeSlots[3]))

		//添加起始时间
		timeslots = append(timeslots, tradingTimeSlots[0])

		//计算每个周期时间段
		for {

			tTmp := timeslots[len(timeslots)-1].Add(duration)

			if !tTmp.After(tradingTimeSlots[1]) {
				timeslots = append(timeslots, tTmp)
				continue
			} else if !tTmp.After(tradingTimeSlots[3]) {

				if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[2]) {
					timeslots = append(timeslots, tTmp)
					continue
				} else {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[0]))
					continue
				}

			} else if !tTmp.After(tradingTimeSlots[5]) {

				if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[4]) {
					timeslots = append(timeslots, tTmp)
					continue
				} else if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[2]) {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[1]))
					continue
				} else {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[1]).Add(nonTradeDuration[0]))
					continue
				}

			} else {
				if !timeslots[len(timeslots)-1].Equal(tradingTimeSlots[5]) {
					timeslots = append(timeslots, tradingTimeSlots[5])
				}
				break
			}
		}

	case ZJ_DAY:

		//交易时间间隔
		nonTradeDuration = append(nonTradeDuration, tradingTimeSlots[2].Sub(tradingTimeSlots[1]))

		//添加起始时间
		timeslots = append(timeslots, tradingTimeSlots[0])

		//计算每个周期时间段
		for {

			tTmp := timeslots[len(timeslots)-1].Add(duration)

			if !tTmp.After(tradingTimeSlots[1]) {
				timeslots = append(timeslots, tTmp)
				continue
			} else if !tTmp.After(tradingTimeSlots[3]) {

				if !timeslots[len(timeslots)-1].Before(tradingTimeSlots[2]) {
					timeslots = append(timeslots, tTmp)
					continue
				} else {
					timeslots = append(timeslots, tTmp.Add(nonTradeDuration[0]))
					continue
				}

			} else {
				if !timeslots[len(timeslots)-1].Equal(tradingTimeSlots[3]) {
					timeslots = append(timeslots, tradingTimeSlots[3])
				}
				break
			}
		}

	default:
		return []time.Time{time.Now()}, Error{"unsupported category"}
	}

	return timeslots, nil

}

//一次返回交易时间段
func GetTradingTimeSlots(instCode string, timeNow time.Time) ([]time.Time, error) {

	var category int //品种的交易时间类型

	switch instCode {
	//大商所品种
	case "a", "b", "bb", "c", "fb", "i", "j", "jd", "jm", "m", "p", "v", "y":
		category = NIGHT_C2
	//大商所品种
	case "cs", "l", "pp":
		category = DAY_C1
	//上期所品种
	case "ag", "au":
		category = NIGHT_C4
	//上期所品种
	case "al", "bu", "cu", "hc", "ni", "pb", "rb", "sn", "zn":
		category = NIGHT_C3
	//上期所品种
	case "ru":
		category = NIGHT_C1
	//上期所品种
	case "fu", "wr":
		category = DAY_C1
	//郑商所品种
	case "CF", "FG", "JR", "LR", "MA", "OI", "PM", "RI", "RM", "RS", "SF", "SM", "SR", "TA", "TC", "WH", "ZC":
		category = NIGHT_C2
	//中金所品种
	case "IC", "IH", "IF", "T", "TF":
		category = ZJ_DAY

	default:
		category = DAY_C1
	}

	//获取每天零时
	//CST时区，UTC向东偏移8小时（28800秒），取零点-8小时
	dayStart := timeNow.In(time.FixedZone("CST", 28800)).Truncate(time.Minute * 60).Add(-time.Hour * time.Duration(timeNow.Hour()))

	//如果是下午3点以后调用该函数，则返回第二天的时间
	if timeNow.Hour() > 18 {
		dayStart = dayStart.Add(time.Hour * 24)
	}

	//交易时间段
	var timeslotsSH []time.Time

	//转换交易时间字符串为时间类型
	dayTradeStartP1, _ := time.Parse(hh_mm_ss, DAY_TRADE_START_P1)
	dayTradeEndP1, _ := time.Parse(hh_mm_ss, DAY_TRADE_END_P1)
	dayTradeStartP2, _ := time.Parse(hh_mm_ss, DAY_TRADE_START_P2)
	dayTradeEndP2, _ := time.Parse(hh_mm_ss, DAY_TRADE_END_P2)
	dayTradeStartP3, _ := time.Parse(hh_mm_ss, DAY_TRADE_START_P3)
	dayTradeEndP3, _ := time.Parse(hh_mm_ss, DAY_TRADE_END_P3)

	switch category {

	case DAY_C1:

		//今天
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP1.Hour())*time.Hour).Add(time.Duration(dayTradeStartP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP1.Hour())*time.Hour).Add(time.Duration(dayTradeEndP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP2.Hour())*time.Hour).Add(time.Duration(dayTradeStartP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP2.Hour())*time.Hour).Add(time.Duration(dayTradeEndP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP3.Hour())*time.Hour).Add(time.Duration(dayTradeStartP3.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP3.Hour())*time.Hour).Add(time.Duration(dayTradeEndP3.Minute())*time.Minute))

		return timeslotsSH, nil

	case NIGHT_C1:

		nightTradeStartP1, _ := time.Parse(hh_mm_ss, NIGHT_C1_TRADE_START_P1)
		nightTradeEndP1, _ := time.Parse(hh_mm_ss, NIGHT_C1_TRADE_END_P1)

		//前一天晚上
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(nightTradeStartP1.Hour()-24)*time.Hour).Add(time.Duration(nightTradeStartP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(nightTradeEndP1.Hour()-24)*time.Hour).Add(time.Duration(nightTradeEndP1.Minute())*time.Minute))

		//今天
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP1.Hour())*time.Hour).Add(time.Duration(dayTradeStartP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP1.Hour())*time.Hour).Add(time.Duration(dayTradeEndP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP2.Hour())*time.Hour).Add(time.Duration(dayTradeStartP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP2.Hour())*time.Hour).Add(time.Duration(dayTradeEndP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP3.Hour())*time.Hour).Add(time.Duration(dayTradeStartP3.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP3.Hour())*time.Hour).Add(time.Duration(dayTradeEndP3.Minute())*time.Minute))

		return timeslotsSH, nil

	case NIGHT_C2:

		nightTradeStartP1, _ := time.Parse(hh_mm_ss, NIGHT_C2_TRADE_START_P1)
		nightTradeEndP1, _ := time.Parse(hh_mm_ss, NIGHT_C2_TRADE_END_P1)

		//前一天晚上
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(nightTradeStartP1.Hour()-24)*time.Hour).Add(time.Duration(nightTradeStartP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(nightTradeEndP1.Hour()-24)*time.Hour).Add(time.Duration(nightTradeEndP1.Minute())*time.Minute))

		//今天
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP1.Hour())*time.Hour).Add(time.Duration(dayTradeStartP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP1.Hour())*time.Hour).Add(time.Duration(dayTradeEndP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP2.Hour())*time.Hour).Add(time.Duration(dayTradeStartP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP2.Hour())*time.Hour).Add(time.Duration(dayTradeEndP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP3.Hour())*time.Hour).Add(time.Duration(dayTradeStartP3.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP3.Hour())*time.Hour).Add(time.Duration(dayTradeEndP3.Minute())*time.Minute))

		return timeslotsSH, nil

	case NIGHT_C3:

		nightTradeStartP1, _ := time.Parse(hh_mm_ss, NIGHT_C3_TRADE_START_P1)
		nightTradeEndP1, _ := time.Parse(hh_mm_ss, NIGHT_C3_TRADE_END_P1)

		//前一天晚上
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(nightTradeStartP1.Hour()-24)*time.Hour).Add(time.Duration(nightTradeStartP1.Minute())*time.Minute))
		//今天
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(nightTradeEndP1.Hour())*time.Hour).Add(time.Duration(nightTradeEndP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP1.Hour())*time.Hour))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP1.Hour())*time.Hour).Add(time.Duration(dayTradeEndP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP2.Hour())*time.Hour).Add(time.Duration(dayTradeStartP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP2.Hour())*time.Hour).Add(time.Duration(dayTradeEndP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP3.Hour())*time.Hour).Add(time.Duration(dayTradeStartP3.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP3.Hour())*time.Hour).Add(time.Duration(dayTradeEndP3.Minute())*time.Minute))

		return timeslotsSH, nil

	case NIGHT_C4:

		nightTradeStartP1, _ := time.Parse(hh_mm_ss, NIGHT_C4_TRADE_START_P1)
		nightTradeEndP1, _ := time.Parse(hh_mm_ss, NIGHT_C4_TRADE_END_P1)

		//前一天晚上
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(nightTradeStartP1.Hour()-24)*time.Hour).Add(time.Duration(nightTradeStartP1.Minute())*time.Minute))

		//今天
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(nightTradeEndP1.Hour())*time.Hour).Add(time.Duration(nightTradeEndP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP1.Hour())*time.Hour).Add(time.Duration(dayTradeStartP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP1.Hour())*time.Hour).Add(time.Duration(dayTradeEndP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP2.Hour())*time.Hour).Add(time.Duration(dayTradeStartP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP2.Hour())*time.Hour).Add(time.Duration(dayTradeEndP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP3.Hour())*time.Hour).Add(time.Duration(dayTradeStartP3.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP3.Hour())*time.Hour).Add(time.Duration(dayTradeEndP3.Minute())*time.Minute))

		return timeslotsSH, nil

	case ZJ_DAY:

		dayTradeStartP1, _ := time.Parse(hh_mm_ss, ZJ_DAY_TRADE_START_P1)
		dayTradeEndP1, _ := time.Parse(hh_mm_ss, ZJ_DAY_TRADE_END_P1)
		dayTradeStartP2, _ := time.Parse(hh_mm_ss, ZJ_DAY_TRADE_START_P2)
		dayTradeEndP2, _ := time.Parse(hh_mm_ss, ZJ_DAY_TRADE_END_P2)

		//今天
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP1.Hour())*time.Hour).Add(time.Duration(dayTradeStartP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP1.Hour())*time.Hour).Add(time.Duration(dayTradeEndP1.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeStartP2.Hour())*time.Hour).Add(time.Duration(dayTradeStartP2.Minute())*time.Minute))
		timeslotsSH = append(timeslotsSH, dayStart.Add(time.Duration(dayTradeEndP2.Hour())*time.Hour).Add(time.Duration(dayTradeEndP2.Minute())*time.Minute))

		return timeslotsSH, nil

	default:
		return []time.Time{time.Now()}, Error{"Invalid instrument category"}
	}
}

//返回Simple K线数据
func GetCandleStickData(kTimeSlots []time.Time, kData []CandleStickData, mktDataS MarketData) []CandleStickData {
	tmpT, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, mktDataS.TimeStamp+" CST")

	var kDataTmp CandleStickData
	var kDataTmpStart time.Time

	//过滤非法行情数据
	if tmpT.Before(kTimeSlots[0]) || tmpT.After(kTimeSlots[len(kTimeSlots)-1]) {
		fmt.Printf("invalid marketdata:{%v},not between %v and %v\n", mktDataS, kTimeSlots[0], kTimeSlots[len(kTimeSlots)-1])
	} else {

		//获取行情数据所在K线时间段
		for i := 0; i < len(kTimeSlots)-1; i++ {

			if !tmpT.Before(kTimeSlots[i]) && !tmpT.After(kTimeSlots[i+1]) {
				kDataTmpStart = kTimeSlots[i]
				break
			}

		}

		kDataTmp.KNum = 1
		kDataTmp.InstrumentID = mktDataS.InstrumentID
		kDataTmp.KTime = kDataTmpStart
		kDataTmp.OpenPrice = mktDataS.LastPrice
		kDataTmp.OpenTimestamp = tmpT
		kDataTmp.ClosePrice = mktDataS.LastPrice
		kDataTmp.CloseTimestamp = tmpT
		kDataTmp.HighestPrice = mktDataS.LastPrice
		kDataTmp.LowestPrice = mktDataS.LastPrice
		kDataTmp.Volume = mktDataS.Volume
		kDataTmp.OpenInterest = mktDataS.OpenInterest
		kDataTmp.Count = 1

		//第一条行情
		if len(kData) == 0 {
			kData = append(kData, kDataTmp)
		} else if kData[len(kData)-1].KTime.Equal(kDataTmpStart) {

			kData[len(kData)-1].ClosePrice = mktDataS.LastPrice
			kData[len(kData)-1].CloseTimestamp = kDataTmp.CloseTimestamp

			if kData[len(kData)-1].HighestPrice < mktDataS.LastPrice {
				kData[len(kData)-1].HighestPrice = mktDataS.LastPrice
			}
			if kData[len(kData)-1].LowestPrice > mktDataS.LastPrice {
				kData[len(kData)-1].LowestPrice = mktDataS.LastPrice
			}
			kData[len(kData)-1].OpenInterest = kDataTmp.OpenInterest
			kData[len(kData)-1].Volume = kDataTmp.Volume
			kData[len(kData)-1].Count++

		} else if kData[len(kData)-1].KTime.Before(kDataTmpStart) {

			kDataTmp.KNum = kData[len(kData)-1].KNum + 1
			kData = append(kData, kDataTmp)

		} else {
			fmt.Println("new market data time is before the lastest candlestick data time")
		}
	}
	return kData
}

//返回Complex K线数据
//包含加权平均价
func GetComplexCandleStickData(kTimeSlots []time.Time, kData []CandleStickData, mktDataS MarketData) []CandleStickData {
	tmpT, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, mktDataS.TimeStamp+" CST")

	var kDataTmp CandleStickData
	var kDataTmpStart time.Time

	//过滤非法行情数据
	if tmpT.Before(kTimeSlots[0]) || tmpT.After(kTimeSlots[len(kTimeSlots)-1]) {
		fmt.Printf("invalid marketdata:{%v},not between %v and %v\n", mktDataS, kTimeSlots[0], kTimeSlots[len(kTimeSlots)-1])
	} else {

		//获取行情数据所在K线时间段
		for i := 0; i < len(kTimeSlots)-1; i++ {

			if !tmpT.Before(kTimeSlots[i]) && !tmpT.After(kTimeSlots[i+1]) {
				kDataTmpStart = kTimeSlots[i]
				break
			}
		}

		kDataTmp.KNum = 1
		kDataTmp.InstrumentID = mktDataS.InstrumentID
		kDataTmp.KTime = kDataTmpStart
		kDataTmp.OpenPrice = mktDataS.LastPrice
		kDataTmp.OpenTimestamp = tmpT
		kDataTmp.ClosePrice = mktDataS.LastPrice
		kDataTmp.CloseTimestamp = tmpT
		kDataTmp.HighestPrice = mktDataS.LastPrice
		kDataTmp.LowestPrice = mktDataS.LastPrice
		kDataTmp.Volume = mktDataS.Volume
		kDataTmp.OpenInterest = mktDataS.OpenInterest
		kDataTmp.Count = 1

		kDataTmp.Expansion.PriceVolumeTotal = mktDataS.LastPrice * float64(mktDataS.Volume)
		kDataTmp.Expansion.VolumeTotal = mktDataS.Volume

		//第一条行情
		if len(kData) == 0 {
			kDataTmp.Expansion.Vwap = kDataTmp.Expansion.PriceVolumeTotal / float64(kDataTmp.Expansion.VolumeTotal)
			kData = append(kData, kDataTmp)
		} else if kData[len(kData)-1].KTime.Equal(kDataTmpStart) {

			kData[len(kData)-1].ClosePrice = mktDataS.LastPrice
			kData[len(kData)-1].CloseTimestamp = kDataTmp.CloseTimestamp

			if kData[len(kData)-1].HighestPrice < mktDataS.LastPrice {
				kData[len(kData)-1].HighestPrice = mktDataS.LastPrice
			}
			if kData[len(kData)-1].LowestPrice > mktDataS.LastPrice {
				kData[len(kData)-1].LowestPrice = mktDataS.LastPrice
			}
			kData[len(kData)-1].OpenInterest = kDataTmp.OpenInterest
			kData[len(kData)-1].Volume = kDataTmp.Volume
			kData[len(kData)-1].Count++

			kData[len(kData)-1].Expansion.PriceVolumeTotal += kDataTmp.Expansion.PriceVolumeTotal
			kData[len(kData)-1].Expansion.VolumeTotal += kDataTmp.Expansion.VolumeTotal
			kData[len(kData)-1].Expansion.Vwap = kData[len(kData)-1].Expansion.PriceVolumeTotal / float64(kData[len(kData)-1].Expansion.VolumeTotal)

		} else if kData[len(kData)-1].KTime.Before(kDataTmpStart) {

			kDataTmp.KNum = kData[len(kData)-1].KNum + 1
			kData = append(kData, kDataTmp)

		} else {
			fmt.Println("new market data time is before the lastest candlestick data time")
		}
	}
	return kData
}
