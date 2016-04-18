package common

import ()

//移动平均线
//支持float64,MarketData和CandleStickData三种类型的数据计算
func MA(data interface{}, length int) (maVal []float64, err error) {

	var sum float64
	var ma []float64

	switch val := data.(type) {

	case []float64:

		if len(val) == 0 {
			return nil, Error{"Input array for function MA() is empty."}
		}

		for i, v := range val {

			sum = 0
			count := 0

			for j := i; j >= 0; j-- {

				if i == 0 {
					sum = v
					count++
				} else if i < length {
					sum = sum + val[j]
					count++
				} else {
					sum = sum + val[j]
					count++
					if i-j >= 4 {
						break
					}
				}
			}

			ma = append(ma, sum/float64(count))
		}
		return ma, nil

	case []MarketData:

		if len(val) == 0 {
			return nil, Error{"Input array for function MA() is empty."}
		}

		for i, v := range val {

			sum = 0
			count := 0

			for j := i; j >= 0; j-- {

				if i == 0 {
					sum = v.LastPrice
					count++
				} else if i < length {
					sum = sum + val[j].LastPrice
					count++
				} else {
					sum = sum + val[j].LastPrice
					count++
					if i-j >= 4 {
						break
					}
				}
			}

			ma = append(ma, sum/float64(count))
		}
		return ma, nil

	case []CandleStickData:

		if len(val) == 0 {
			return nil, Error{"Input array for function MA() is empty."}
		}

		for i, v := range val {

			sum = 0
			count := 0

			for j := i; j >= 0; j-- {

				if i == 0 {
					sum = v.ClosePrice
					count++
				} else if i < length {
					sum = sum + val[j].ClosePrice
					count++
				} else {
					sum = sum + val[j].ClosePrice
					count++
					if i-j >= 4 {
						break
					}
				}
			}

			ma = append(ma, sum/float64(count))
		}
		return ma, nil

	default:
		return nil, Error{"Input type is not supported by function MA()."}
	}

}

//简单移动平均
//支持float64,MarketData和CandleStickData三种类型的数据计算
//SMA=-1，视为无效数据
func SMA(data interface{}, length int) (smaVal []float64, err error) {

	var sma []float64
	var smaTmp float64

	switch val := data.(type) {

	case []float64:

		if len(val) < length {
			return nil, Error{"Input array for function SMA() is less than length."}
		}

		for i, _ := range val {
			if i < length-1 {
				continue
			} else {
				smaTmp, _ = Avg(val[i-length+1 : i+1])
				sma = append(sma, smaTmp)
			}
		}

		return sma, nil

	case []MarketData:

		if len(val) < length {
			return nil, Error{"Input array for function SMA() is less than length."}
		}

		var valTmp []float64

		for _, v := range val {
			valTmp = append(valTmp, v.LastPrice)
		}

		for i, _ := range val {
			if i < length-1 {
				continue
			} else {
				smaTmp, _ = Avg(valTmp[i-length+1 : i+1])
				sma = append(sma, smaTmp)
			}
		}

		return sma, nil

	case []CandleStickData:

		if len(val) <= length {
			return nil, Error{"Input array for function SMA() is less than length."}
		}

		var valTmp []float64

		for _, v := range val {
			valTmp = append(valTmp, v.ClosePrice)
		}

		for i, _ := range valTmp {
			if i < length-1 {
				continue
			} else {
				smaTmp, _ = Avg(valTmp[i-length+1 : i+1])
				sma = append(sma, smaTmp)
			}
		}

		return sma, nil

	default:
		return nil, Error{"Input type is not supported by function SMA()."}
	}
}

//指数移动平均
//支持float64,MarketData和CandleStickData三种类型的数据计算
//EMAtoday=α * Pricetoday + ( 1 - α ) * EMAyesterday; α=2/(length+1)
func EMA(data interface{}, length int) (emaVal []float64, err error) {

	var alpha float64 = float64(2) / float64(length+1) //平滑指数
	var ema []float64

	switch val := data.(type) {

	case []float64:

		if len(val) == 0 {
			return nil, Error{"Input array for function EMA() is empty."}
		}

		if len(val) == 1 {
			ema = append(ema, val[0])
		} else {
			for i, _ := range val {
				if i == 0 {
					ema = append(ema, val[0])
				} else {
					ema = append(ema, alpha*val[i]+(1-alpha)*ema[i-1])
				}
			}
		}

		return ema, nil

	case []MarketData:

		if len(val) == 0 {
			return nil, Error{"Input array for function EMA() is empty."}
		}

		if len(val) == 1 {
			ema = append(ema, val[0].LastPrice)
		} else {
			for i, _ := range val {
				if i == 0 {
					ema = append(ema, val[0].LastPrice)
				} else {
					ema = append(ema, alpha*val[i].LastPrice+(1-alpha)*ema[i-1])
				}
			}
		}

		return ema, nil

	case []CandleStickData:

		if len(val) == 0 {
			return nil, Error{"Input array for function EMA() is empty."}
		}

		if len(val) == 1 {
			ema = append(ema, val[0].ClosePrice)
		} else {
			for i, _ := range val {
				if i == 0 {
					ema = append(ema, val[0].ClosePrice)
				} else {
					ema = append(ema, alpha*val[i].ClosePrice+(1-alpha)*ema[i-1])
				}
			}
		}

		return ema, nil

	default:
		return nil, Error{"Input type is not supported by function EMA()."}
	}
}

//平滑异同移动平均线
//DIFF = EMA(close,12) - EMA(close,26)
//DEA
func MACD(data interface{}, fastLen, slowLen, macdLen int) (macdVal []MACDData, err error) {
	var macd []MACDData
	var macdTmp MACDData
	var dif []float64
	var dea []float64

	emaFast, _ := EMA(data, fastLen)
	emaSlow, _ := EMA(data, slowLen)

	for i, _ := range emaFast {
		dif = append(dif, emaFast[i]-emaSlow[i])
	}

	emaDif, _ := EMA(dif, macdLen)

	for _, v := range dif {
		dea = append(dea, v)
	}

	for j, _ := range emaFast {
		macdTmp.EMAFast = emaFast[j]
		macdTmp.EMASlow = emaSlow[j]
		macdTmp.Dif = emaDif[j]
		macdTmp.Dea = dea[j]
		macd = append(macd, macdTmp)
	}

	return macd, nil
}

//相对强弱指数
//RSI = SMA(U,n)/(SMA(U,n)+SMA(D,n)) * 100; U价格上涨，D价格下跌
func RSI(data interface{}, length int) (rsiVal []float64, err error) {
	var rsi []float64

	var smaU []float64
	var smaD []float64

	switch val := data.(type) {

	case []float64:

		for i, _ := range val {
			if i < length-1 {
				rsi = append(rsi, -1)
			} else {
				var dataU []float64
				var dataD []float64

				for j := i - length + 1; j <= i; j++ {

					if j == 0 {
						rsi = append(rsi, -1)
					} else if val[j] >= val[j-1] {
						dataU = append(dataU, val[j]-val[j-1])
						dataD = append(dataD, 0)
					} else {
						dataU = append(dataU, 0)
						dataD = append(dataD, val[j-1]-val[j])
					}
				}

				smaUTmp, err1 := SMA(dataU, length)
				smaDTmp, err2 := SMA(dataD, length)

				if err1 == nil && err2 == nil {
					smaU = append(smaU, smaUTmp[0])
					smaD = append(smaD, smaDTmp[0])
				}
			}
		}

		for m, _ := range smaU {
			if smaU[m] == 0 && smaD[m] == 0 {
				rsi = append(rsi, 50)
			} else {
				rsi = append(rsi, 100*smaU[m]/(smaU[m]+smaD[m]))
			}
		}

		return rsi, nil

	case []MarketData:

		for i, _ := range val {
			if i < length-1 {
				rsi = append(rsi, -1)
			} else {
				var dataU []float64
				var dataD []float64

				for j := i - length + 1; j <= i; j++ {

					if j == 0 {
						rsi = append(rsi, -1)
					} else if val[j].LastPrice >= val[j-1].LastPrice {
						dataU = append(dataU, val[j].LastPrice-val[j-1].LastPrice)
						dataD = append(dataD, 0)
					} else {
						dataU = append(dataU, 0)
						dataD = append(dataD, val[j-1].LastPrice-val[j].LastPrice)
					}
				}

				smaUTmp, err1 := SMA(dataU, length)
				smaDTmp, err2 := SMA(dataD, length)

				if err1 == nil && err2 == nil {
					smaU = append(smaU, smaUTmp[0])
					smaD = append(smaD, smaDTmp[0])
				}
			}
		}

		for m, _ := range smaU {
			if smaU[m] == 0 && smaD[m] == 0 {
				rsi = append(rsi, 50)
			} else {
				rsi = append(rsi, 100*smaU[m]/(smaU[m]+smaD[m]))
			}
		}

		return rsi, nil

	case []CandleStickData:

		for i, _ := range val {
			if i < length-1 {
				rsi = append(rsi, -1)
			} else {
				var dataU []float64
				var dataD []float64

				for j := i - length + 1; j <= i; j++ {

					if j == 0 {
						rsi = append(rsi, -1)
					} else if val[j].ClosePrice >= val[j-1].ClosePrice {
						dataU = append(dataU, val[j].ClosePrice-val[j-1].ClosePrice)
						dataD = append(dataD, 0)
					} else {
						dataU = append(dataU, 0)
						dataD = append(dataD, val[j-1].ClosePrice-val[j].ClosePrice)
					}
				}

				smaUTmp, err1 := SMA(dataU, length)
				smaDTmp, err2 := SMA(dataD, length)

				if err1 == nil && err2 == nil {
					smaU = append(smaU, smaUTmp[0])
					smaD = append(smaD, smaDTmp[0])
				}
			}
		}

		for m, _ := range smaU {
			if smaU[m] == 0 && smaD[m] == 0 {
				rsi = append(rsi, 50)
			} else {
				rsi = append(rsi, 100*smaU[m]/(smaU[m]+smaD[m]))
			}
		}

		return rsi, nil

	default:
		return rsi, nil
	}
}

//未成熟随机值（Raw Stochastic Value)
//RSV=(Cn-Ln)/(Hn-Ln)*100%
//n：是经过的交易期间（一般订为9，n=9）；
//Cn：是第n个周期的收盘价；
//Hn,Ln：分别是过去n周期内的最高价和最低价，一般以9为基准。
func RSV(data interface{}, length int) (rsvVal []float64, err error) {

	var rsv []float64
	var rsvTmp float64

	switch val := data.(type) {

	case []float64:

		if len(val) == 0 {
			return nil, Error{"Input array for function RSV() is empty."}
		}

		for i, _ := range val {

			if i < length {

				_, High, _ := Max(val[0 : i+1])
				_, Low, _ := Min(val[0 : i+1])
				if High == Low {
					rsvTmp = 0
				} else {
					rsvTmp = (val[i] - Low.(float64)) / (High.(float64) - Low.(float64))
				}

			} else {
				_, High, _ := Max(val[i-length+1 : i+1])
				_, Low, _ := Min(val[i-length+1 : i+1])
				rsvTmp = (val[i] - Low.(float64)) / (High.(float64) - Low.(float64))
			}
			rsv = append(rsv, rsvTmp*100)
		}

		return rsv, nil

	case []MarketData:

		if len(val) == 0 {
			return nil, Error{"Input array for function RSV() is empty."}
		}

		var valTmp []float64

		for _, v := range val {
			valTmp = append(valTmp, v.LastPrice)
		}

		for i, _ := range valTmp {

			if i < length {

				_, High, _ := Max(valTmp[0 : i+1])
				_, Low, _ := Min(valTmp[0 : i+1])
				if High == Low {
					rsvTmp = 0
				} else {
					rsvTmp = (valTmp[i] - Low.(float64)) / (High.(float64) - Low.(float64))
				}

			} else {
				_, High, _ := Max(valTmp[i-length+1 : i+1])
				_, Low, _ := Min(valTmp[i-length+1 : i+1])
				rsvTmp = (valTmp[i] - Low.(float64)) / (High.(float64) - Low.(float64))
			}
			rsv = append(rsv, rsvTmp*100)
		}

		return rsv, nil

	case []CandleStickData:

		if len(val) == 0 {
			return nil, Error{"Input array for function RSV() is empty."}
		}

		var valTmpClose []float64
		var valTmpHigh []float64
		var valTmpLow []float64

		for _, v := range val {
			valTmpClose = append(valTmpClose, v.ClosePrice)
			valTmpHigh = append(valTmpHigh, v.HighestPrice)
			valTmpLow = append(valTmpLow, v.LowestPrice)
		}

		for i, _ := range valTmpClose {

			if i < length {

				_, High, _ := Max(valTmpHigh[0 : i+1])
				_, Low, _ := Min(valTmpLow[0 : i+1])
				if High == Low {
					rsvTmp = 0
				} else {
					rsvTmp = (valTmpClose[i] - Low.(float64)) / (High.(float64) - Low.(float64))
				}

			} else {
				_, High, _ := Max(valTmpHigh[i-length+1 : i+1])
				_, Low, _ := Min(valTmpLow[i-length+1 : i+1])
				rsvTmp = (valTmpClose[i] - Low.(float64)) / (High.(float64) - Low.(float64))
			}
			rsv = append(rsv, rsvTmp*100)
		}

		return rsv, nil

	default:
		return nil, Error{"Input type is not supported by function RSV()."}
	}
}

//KD随机指标
//K=EMA(RSV,N)
//D=EMA(K,N)
func KD(data interface{}, length int) (KDVal []KDData, err error) {

	var kd []KDData

	rsv, err := RSV(data, length)
	if err != nil {
		return kd, err
	}

	k, err1 := EMA(rsv, length)
	if err1 != nil {
		return kd, err1
	}

	d, err2 := EMA(k, length)
	if err2 != nil {
		return kd, err2
	}

	for i, _ := range k {
		kd = append(kd, KDData{k[i], d[i]})
	}
	return kd, nil

}
