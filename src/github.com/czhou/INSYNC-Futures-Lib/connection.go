package common

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"sync"
	"time"
)

//行情数据处理
type ProcessMarketData interface {
	ProcessMarketData(marketData MarketData) (interface{}, error)
	ProcessCandleStickData(candleStickData []CandleStickData) (interface{}, error)
}

//持仓数据处理
type ProcessPositionData interface {
	ProcessPositionData(posData []PositionData) (interface{}, error)
}

//下单响应处理
type ProcessPlaceOrderResponseData interface {
	ProcessPlaceOrderResponseData(porData []PlaceOrderResponseData) (interface{}, error)
}

//成交回报处理
type ProcessOrderMatchedReturnData interface {
	ProcessOrderMatchedReturnData(omrData []OrderMatchedReturnData) (interface{}, error)
}

//撤单响应处理
type ProcessCancelOrderResponseData interface {
	ProcessCancelOrderResponseData(pcoData []CancelOrderResponseData) (interface{}, error)
}

//撤单回报处理
type ProcessOrderCanceledReturnData interface {
	ProcessOrderCanceledReturnData(pocData []OrderCanceledReturnData) (interface{}, error)
}

//资金数据处理
type ProcessCapitalData interface {
	ProcessCapitalData(cptData CapitalData) (interface{}, error)
}

type Error struct {
	err string
}

func (e Error) Error() string { return e.err }

type RedisConn struct {
	redis.Conn
	connStr string
}

type RedisSubConn struct {
	SubConn        *redis.PubSubConn
	ChannelName    string //频道名称
	InstrumentCode string //合约代码
	InstrumentNum  string //合约期号
	IsUnnatual     bool   //是否为非自然行情，之所以没用IsNatual，是因为需要默认值false
}

//创建连接
func NewConn(hostNport string) (*RedisConn, error) {
	conn, err := redis.Dial("tcp", hostNport)
	if conn == nil {
		return nil, err
	}

	return &RedisConn{conn, hostNport}, nil
}

//错误返回
func (c *RedisConn) Error(err string) string {
	return err
}

//返回连接串
func (c *RedisConn) GetConnStr() string {
	return c.connStr
}

//订阅MarketData频道：SubscribeChan("DB8:futures:MarketData", "SR", "601", true)
//订阅Query频道：SubscribeChan("DB8:futures:ReqQryPosition", "", "", false)
func (c *RedisConn) SubscribeChan(chanName, instCode, instNum string, isMarketData bool) (RedisSubConn, error) {
	redisSubCon := RedisSubConn{}

	channel := "QUOTE:" + instCode + instNum

	if isMarketData {
		//通知reidis需要订阅的品种
		c.Do("PUBLISH", chanName, "{\"instrumentCode\":\""+instCode+"\",\"instrumentNo\":\""+instNum+"\"}")
		redisSubCon.ChannelName = chanName
		redisSubCon.InstrumentCode = instCode
		redisSubCon.InstrumentNum = instNum
	} else {
		channel = chanName
		redisSubCon.ChannelName = channel
	}

	//订阅频道
	subConn := redis.PubSubConn{c}

	if err := subConn.Subscribe(channel); err != nil {
		return redisSubCon, err
	}

	redisSubCon.SubConn = &subConn

	return redisSubCon, nil

}

//填充MarketData数据和K线数据
//kDuration:K线周期
/*例:
逐笔行情:kDuration=0;
2分钟K线:kDuration=time.Minute*2;
2小时K线:kDuration=time.Hour*2;
*/
func (c *RedisConn) LoadMarketData(subConn RedisSubConn, process ProcessMarketData, kParams []KParams) error {

	//监听行情并填充行情数据
	if subConn.SubConn == nil {
		return Error{"nill pointer of PubSubConn."}
	}

	var mktDataS []MarketData          //行情数据
	var currentTime time.Time          //当前时间
	sysT := GetSysTime()               //获取回测系统时间单例
	var isKTimeSlotsReady bool = false //K线槽是否已经生成

	kData := make(map[time.Duration][]CandleStickData) //不同周期的K线数据
	kTimeSlots := make(map[time.Duration][]time.Time)  //不同周期的时间切片

	if len(kParams) == 0 {
		kParams = append(kParams, KParams{K_NONE, 0})
	}

	var isLastKDataProcessed bool = false      //是否已生成最后一根K线
	var isLastKDataRoutineCreated bool = false //自然行情下，生成最后一根K线的routine是否已开始

	for {
		switch v := subConn.SubConn.Receive().(type) {

		case redis.Message:

			//处理JSON数据
			data := string(v.Data)
			data = strings.Replace(data, "'", "\"", -1)
			js, err := simplejson.NewJson([]byte(data))

			if err != nil {
				fmt.Println(err)
				continue
			}

			if !subConn.IsUnnatual {
				//自然行情
				mktDataS, currentTime = processMarketDataJS(js, !subConn.IsUnnatual)
				//fmt.Println("T[1]:", mktDataS)
			} else {
				mktDataS, currentTime = processMarketDataJS(js, !subConn.IsUnnatual)
				//fmt.Println("T[2]:", mktDataS)
			}

			//同步回测系统时间
			sysT.sync(currentTime)
			//fmt.Println("Natual Time:", sysT.Now(true))
			//fmt.Println("Unnatual Time:", sysT.Now(false))

			//非回测情况下,14:59分的时候创建一个1分钟进程，用来完成15:00最后一根K线的生成
			if !subConn.IsUnnatual && !isLastKDataProcessed && !isLastKDataRoutineCreated && currentTime.Hour() == 14 && currentTime.Minute() >= 59 {

				go func() {
					fmt.Println("Processing last K bar routine is started.")
					time.Sleep(time.Minute)
					for _, v := range kParams {
						if v.KComplexity == K_NONE {
							continue
						}
						process.ProcessCandleStickData(kData[v.KDuration])
						isLastKDataProcessed = true
					}
					fmt.Printf("Last K bar of the day[%v] is generated. \n", currentTime)
				}()

				isLastKDataRoutineCreated = true
			}

			//换天，生成下一天的K线槽
			if len(kTimeSlots) == 0 {
				isKTimeSlotsReady = false
			} else if currentTime.After(kTimeSlots[kParams[0].KDuration][len(kTimeSlots[kParams[0].KDuration])-1]) {
				isKTimeSlotsReady = false
			}

			//生成K线槽
			if !isKTimeSlotsReady {
				for _, v := range kParams {
					//获取一天的K线时间切片
					var kTimeSlotsTmp []time.Time
					if !subConn.IsUnnatual {
						kTimeSlotsTmp, err = GetIntervalTimeSlots(subConn.InstrumentCode, v.KDuration, time.Now())
					} else {
						kTimeSlotsTmp, err = GetIntervalTimeSlots(subConn.InstrumentCode, v.KDuration, currentTime)
					}
					if err != nil {
						fmt.Printf("GetIntervalTimeSlots(%v,%v) failed, in func LoadMarketData\n", subConn.InstrumentCode, v.KDuration)
						continue
					}

					kTimeSlots[v.KDuration] = kTimeSlotsTmp
				}

				isKTimeSlotsReady = true
			}

			//分时数据处理
			for _, v := range kParams {

				//回测情况下，在15:00生成最后一根K线
				if subConn.IsUnnatual {
					if currentTime.Hour() == 15 && !isLastKDataProcessed {
						process.ProcessCandleStickData(kData[v.KDuration])
						isLastKDataProcessed = true
						fmt.Printf("Last K bar of the day[%v] is generated. \n", currentTime)
					}
				}

				for _, mkt := range mktDataS {

					switch v.KComplexity {
					case K_NONE:

						//过滤无效行情数据
						if currentTime.Before(kTimeSlots[v.KDuration][0]) || currentTime.After(kTimeSlots[v.KDuration][len(kTimeSlots[v.KDuration])-1]) {

							fmt.Printf("Invalid MarketData current time[%v] not in[%v,%v]:%v\n", currentTime, kTimeSlots[v.KDuration][0], kTimeSlots[v.KDuration][len(kTimeSlots[v.KDuration])-1], mkt)
						} else {

							//fmt.Println("[MarketData]:", mkt)
							process.ProcessMarketData(mkt)
						}

					case K_SIMPLE:

						//过滤无效行情数据
						if currentTime.Before(kTimeSlots[v.KDuration][0]) || currentTime.After(kTimeSlots[v.KDuration][len(kTimeSlots[v.KDuration])-1]) {
							fmt.Println("Invalid MarketData:", mkt)
						} else {
							kData[v.KDuration] = GetCandleStickData(kTimeSlots[v.KDuration], kData[v.KDuration], mkt)
							if len(kData[v.KDuration]) > 0 {
								kData[v.KDuration][len(kData[v.KDuration])-1].KDuration = v.KDuration

								if len(kData[v.KDuration]) > 1 {
									//K线下一周期的第一条行情到达的时候触发处理上一条完整地K线
									if kData[v.KDuration][len(kData[v.KDuration])-1].KNum == kData[v.KDuration][len(kData[v.KDuration])-2].KNum+1 && kData[v.KDuration][len(kData[v.KDuration])-1].Count == 1 {
										//每天最后一根K线不由下一条K线行情来触发
										if !isLastKDataProcessed {
											//只返回完整的K线
											process.ProcessCandleStickData(kData[v.KDuration][0 : len(kData[v.KDuration])-1])
											isLastKDataProcessed = false
										}
									}
								}
							}
						}

					case K_COMPLEX:

						//过滤无效行情数据
						if currentTime.Before(kTimeSlots[v.KDuration][0]) || currentTime.After(kTimeSlots[v.KDuration][len(kTimeSlots[v.KDuration])-1]) {
							fmt.Println("Invalid MarketData:", mkt)
						} else {
							kData[v.KDuration] = GetComplexCandleStickData(kTimeSlots[v.KDuration], kData[v.KDuration], mkt)
							if len(kData[v.KDuration]) > 0 {
								kData[v.KDuration][len(kData[v.KDuration])-1].KDuration = v.KDuration

								if len(kData[v.KDuration]) > 1 {
									//K线下一周期的第一条行情到达的时候触发处理上一条完整地K线
									if kData[v.KDuration][len(kData[v.KDuration])-1].KNum == kData[v.KDuration][len(kData[v.KDuration])-2].KNum+1 && kData[v.KDuration][len(kData[v.KDuration])-1].Count == 1 {
										//每天最后一根K线不由下一条K线行情来触发
										if !isLastKDataProcessed {
											//只返回完整的K线
											process.ProcessCandleStickData(kData[v.KDuration][0 : len(kData[v.KDuration])-1])
											isLastKDataProcessed = false
										}
									}
								}
							}
						}

					default:
					}
				}
			}

		case redis.Subscription:
			fmt.Println(v.Channel, v.Count, v.Kind)
		case error:
			return v
		}
	}

}

//处理行情JSON,分两种情况：自然播放和回放
func processMarketDataJS(js *simplejson.Json, isNatual bool) ([]MarketData, time.Time) {

	var mktDataS []MarketData
	var tmpT time.Time

	if isNatual {
		mktTmp := MarketData{}

		//填充模型数据
		tradingDay, _ := js.Get("tradingDay").String()
		instrumentID, _ := js.Get("instrumentID").String()
		exchangeID, _ := js.Get("exchangeID").String()
		exchangeInstID, _ := js.Get("exchangeInstID").String()
		lastPrice, _ := js.Get("lastPrice").String()
		preSettlementPrice, _ := js.Get("preSettlementPrice").String()
		preClosePrice, _ := js.Get("preClosePrice").String()
		preOpenInterest, _ := js.Get("preOpenInterest").String()
		openPrice, _ := js.Get("openPrice").String()
		highestPrice, _ := js.Get("highestPrice").String()
		lowestPrice, _ := js.Get("lowestPrice").String()
		volume, _ := js.Get("Volume").String()
		turnover, _ := js.Get("turnover").String()
		openInterest, _ := js.Get("openInterest").String()
		closePrice, _ := js.Get("closePrice").String()
		settlementPrice, _ := js.Get("settlementPrice").String()
		upperLimitPrice, _ := js.Get("upperLimitPrice").String()
		lowerLimitPrice, _ := js.Get("lowerLimitPrice").String()
		preDelta, _ := js.Get("preDelta").String()
		currentDelta, _ := js.Get("currDelta").String()
		updateTime, _ := js.Get("UpdateTime").String()
		updateMillisec, _ := js.Get("UpdateMillisec").String()
		bidPrice1, _ := js.Get("BidPrice1").String()
		bidVolume1, _ := js.Get("BidVolume1").String()
		askPrice1, _ := js.Get("AskPrice1").String()
		askVolume1, _ := js.Get("AskVolume1").String()
		bidPrice2, _ := js.Get("BidPrice2").String()
		bidVolume2, _ := js.Get("BidVolume2").String()
		askPrice2, _ := js.Get("AskPrice2").String()
		askVolume2, _ := js.Get("AskVolume2").String()
		bidPrice3, _ := js.Get("BidPrice3").String()
		bidVolume3, _ := js.Get("BidVolume3").String()
		askPrice3, _ := js.Get("AskPrice3").String()
		askVolume3, _ := js.Get("AskVolume3").String()
		bidPrice4, _ := js.Get("BidPrice4").String()
		bidVolume4, _ := js.Get("BidVolume4").String()
		askPrice4, _ := js.Get("AskPrice4").String()
		askVolume4, _ := js.Get("AskVolume4").String()
		bidPrice5, _ := js.Get("BidPrice5").String()
		bidVolume5, _ := js.Get("BidVolume5").String()
		askPrice5, _ := js.Get("AskPrice5").String()
		askVolume5, _ := js.Get("AskVolume5").String()
		averagePrice, _ := js.Get("AveragePrice").String()
		xpeedTime, _ := js.Get("XSpeedTime").String()
		stdTime, _ := js.Get("stdTime").String()

		mktTmp.TradingDate = tradingDay
		mktTmp.InstrumentID = instrumentID
		mktTmp.ExchangeID = exchangeID
		mktTmp.ExchangeInstID = exchangeInstID
		mktTmp.LastPrice, _ = strconv.ParseFloat(lastPrice, 64)
		mktTmp.PreSettlementPrice, _ = strconv.ParseFloat(preSettlementPrice, 64)
		mktTmp.PreClosePrice, _ = strconv.ParseFloat(preClosePrice, 64)
		mktTmp.PreOpenInterest, _ = strconv.ParseInt(preOpenInterest, 0, 64)
		mktTmp.OpenPrice, _ = strconv.ParseFloat(openPrice, 64)
		mktTmp.HighestPrice, _ = strconv.ParseFloat(highestPrice, 64)
		mktTmp.LowestPrice, _ = strconv.ParseFloat(lowestPrice, 64)
		mktTmp.Volume, _ = strconv.ParseInt(volume, 0, 64)
		mktTmp.Turnover, _ = strconv.ParseFloat(turnover, 64)
		mktTmp.OpenInterest, _ = strconv.ParseInt(openInterest, 0, 64)
		mktTmp.ClosePrice, _ = strconv.ParseFloat(closePrice, 64)
		mktTmp.SettlementPrice, _ = strconv.ParseFloat(settlementPrice, 64)
		mktTmp.UpperLimitPrice, _ = strconv.ParseFloat(upperLimitPrice, 64)
		mktTmp.LowerLimitPrice, _ = strconv.ParseFloat(lowerLimitPrice, 64)
		mktTmp.PreDelta, _ = strconv.ParseFloat(preDelta, 64)
		mktTmp.CurrentDelta, _ = strconv.ParseFloat(currentDelta, 64)
		mktTmp.UpdateTime = updateTime
		mktTmp.UpdateMillisec = updateMillisec
		mktTmp.BidPrice1, _ = strconv.ParseFloat(bidPrice1, 64)
		mktTmp.BidVolume1, _ = strconv.ParseInt(bidVolume1, 0, 64)
		mktTmp.AskPrice1, _ = strconv.ParseFloat(askPrice1, 64)
		mktTmp.AskVolume1, _ = strconv.ParseInt(askVolume1, 0, 64)
		mktTmp.BidPrice2, _ = strconv.ParseFloat(bidPrice2, 64)
		mktTmp.BidVolume2, _ = strconv.ParseInt(bidVolume2, 0, 64)
		mktTmp.AskPrice2, _ = strconv.ParseFloat(askPrice2, 64)
		mktTmp.AskVolume2, _ = strconv.ParseInt(askVolume2, 0, 64)
		mktTmp.BidPrice3, _ = strconv.ParseFloat(bidPrice3, 64)
		mktTmp.BidVolume3, _ = strconv.ParseInt(bidVolume3, 0, 64)
		mktTmp.AskPrice3, _ = strconv.ParseFloat(askPrice3, 64)
		mktTmp.AskVolume3, _ = strconv.ParseInt(askVolume3, 0, 64)
		mktTmp.BidPrice4, _ = strconv.ParseFloat(bidPrice4, 64)
		mktTmp.BidVolume4, _ = strconv.ParseInt(bidVolume4, 0, 64)
		mktTmp.AskPrice4, _ = strconv.ParseFloat(askPrice4, 64)
		mktTmp.AskVolume4, _ = strconv.ParseInt(askVolume4, 0, 64)
		mktTmp.BidPrice5, _ = strconv.ParseFloat(bidPrice5, 64)
		mktTmp.BidVolume5, _ = strconv.ParseInt(bidVolume5, 0, 64)
		mktTmp.AskPrice5, _ = strconv.ParseFloat(askPrice5, 64)
		mktTmp.AskVolume5, _ = strconv.ParseInt(askVolume5, 0, 64)
		mktTmp.AveragePrice, _ = strconv.ParseFloat(averagePrice, 64)
		mktTmp.XSpeedTime = xpeedTime
		mktTmp.TimeStamp = stdTime

		tmpT, _ = time.Parse(YYYY_MM_DD_hh_mm_ss_MST, mktTmp.TimeStamp+" CST")
		mktDataS = append(mktDataS, mktTmp)

		//fmt.Println("IsNatual TmpT", tmpT)

		//return mktDataS, time.Now()
		return mktDataS, tmpT

	} else {
		mktTmp := MarketData{}

		time_now, _ := js.Get("time_now").String()
		ticks, _ := js.Get("ticks").Array()

		for _, v := range ticks {
			data, _ := v.(map[string]interface{})

			//填充模型数据
			tradingDay := data["tradingDay"].(string)
			instrumentID := data["instrumentID"].(string)
			exchangeID := data["exchangeID"].(string)
			exchangeInstID := data["exchangeInstID"].(string)
			lastPrice := data["lastPrice"].(string)
			preSettlementPrice := data["preSettlementPrice"].(string)
			preClosePrice := data["preClosePrice"].(string)
			preOpenInterest := data["preOpenInterest"].(string)
			openPrice := data["openPrice"].(string)
			highestPrice := data["highestPrice"].(string)
			lowestPrice := data["lowestPrice"].(string)
			volume := data["Volume"].(string)
			turnover := data["turnover"].(string)
			openInterest := data["openInterest"].(string)
			closePrice := data["closePrice"].(string)
			settlementPrice := data["settlementPrice"].(string)
			upperLimitPrice := data["upperLimitPrice"].(string)
			lowerLimitPrice := data["lowerLimitPrice"].(string)
			preDelta := data["preDelta"].(string)
			currentDelta := data["currDelta"].(string)
			updateTime := data["UpdateTime"].(string)
			updateMillisec := data["UpdateMillisec"].(string)
			bidPrice1 := data["BidPrice1"].(string)
			bidVolume1 := data["BidVolume1"].(string)
			askPrice1 := data["AskPrice1"].(string)
			askVolume1 := data["AskVolume1"].(string)
			bidPrice2 := data["BidPrice2"].(string)
			bidVolume2 := data["BidVolume2"].(string)
			askPrice2 := data["AskPrice2"].(string)
			askVolume2 := data["AskVolume2"].(string)
			bidPrice3 := data["BidPrice3"].(string)
			bidVolume3 := data["BidVolume3"].(string)
			askPrice3 := data["AskPrice3"].(string)
			askVolume3 := data["AskVolume3"].(string)
			bidPrice4 := data["BidPrice4"].(string)
			bidVolume4 := data["BidVolume4"].(string)
			askPrice4 := data["AskPrice4"].(string)
			askVolume4 := data["AskVolume4"].(string)
			bidPrice5 := data["BidPrice5"].(string)
			bidVolume5 := data["BidVolume5"].(string)
			askPrice5 := data["AskPrice5"].(string)
			askVolume5 := data["AskVolume5"].(string)
			averagePrice := data["AveragePrice"].(string)
			xpeedTime := data["XSpeedTime"].(string)
			stdTime := data["stdTime"].(string)

			mktTmp.TradingDate = tradingDay
			mktTmp.InstrumentID = instrumentID
			mktTmp.ExchangeID = exchangeID
			mktTmp.ExchangeInstID = exchangeInstID
			mktTmp.LastPrice, _ = strconv.ParseFloat(lastPrice, 64)
			mktTmp.PreSettlementPrice, _ = strconv.ParseFloat(preSettlementPrice, 64)
			mktTmp.PreClosePrice, _ = strconv.ParseFloat(preClosePrice, 64)
			mktTmp.PreOpenInterest, _ = strconv.ParseInt(preOpenInterest, 0, 64)
			mktTmp.OpenPrice, _ = strconv.ParseFloat(openPrice, 64)
			mktTmp.HighestPrice, _ = strconv.ParseFloat(highestPrice, 64)
			mktTmp.LowestPrice, _ = strconv.ParseFloat(lowestPrice, 64)
			mktTmp.Volume, _ = strconv.ParseInt(volume, 0, 64)
			mktTmp.Turnover, _ = strconv.ParseFloat(turnover, 64)
			mktTmp.OpenInterest, _ = strconv.ParseInt(openInterest, 0, 64)
			mktTmp.ClosePrice, _ = strconv.ParseFloat(closePrice, 64)
			mktTmp.SettlementPrice, _ = strconv.ParseFloat(settlementPrice, 64)
			mktTmp.UpperLimitPrice, _ = strconv.ParseFloat(upperLimitPrice, 64)
			mktTmp.LowerLimitPrice, _ = strconv.ParseFloat(lowerLimitPrice, 64)
			mktTmp.PreDelta, _ = strconv.ParseFloat(preDelta, 64)
			mktTmp.CurrentDelta, _ = strconv.ParseFloat(currentDelta, 64)
			mktTmp.UpdateTime = updateTime
			mktTmp.UpdateMillisec = updateMillisec
			mktTmp.BidPrice1, _ = strconv.ParseFloat(bidPrice1, 64)
			mktTmp.BidVolume1, _ = strconv.ParseInt(bidVolume1, 0, 64)
			mktTmp.AskPrice1, _ = strconv.ParseFloat(askPrice1, 64)
			mktTmp.AskVolume1, _ = strconv.ParseInt(askVolume1, 0, 64)
			mktTmp.BidPrice2, _ = strconv.ParseFloat(bidPrice2, 64)
			mktTmp.BidVolume2, _ = strconv.ParseInt(bidVolume2, 0, 64)
			mktTmp.AskPrice2, _ = strconv.ParseFloat(askPrice2, 64)
			mktTmp.AskVolume2, _ = strconv.ParseInt(askVolume2, 0, 64)
			mktTmp.BidPrice3, _ = strconv.ParseFloat(bidPrice3, 64)
			mktTmp.BidVolume3, _ = strconv.ParseInt(bidVolume3, 0, 64)
			mktTmp.AskPrice3, _ = strconv.ParseFloat(askPrice3, 64)
			mktTmp.AskVolume3, _ = strconv.ParseInt(askVolume3, 0, 64)
			mktTmp.BidPrice4, _ = strconv.ParseFloat(bidPrice4, 64)
			mktTmp.BidVolume4, _ = strconv.ParseInt(bidVolume4, 0, 64)
			mktTmp.AskPrice4, _ = strconv.ParseFloat(askPrice4, 64)
			mktTmp.AskVolume4, _ = strconv.ParseInt(askVolume4, 0, 64)
			mktTmp.BidPrice5, _ = strconv.ParseFloat(bidPrice5, 64)
			mktTmp.BidVolume5, _ = strconv.ParseInt(bidVolume5, 0, 64)
			mktTmp.AskPrice5, _ = strconv.ParseFloat(askPrice5, 64)
			mktTmp.AskVolume5, _ = strconv.ParseInt(askVolume5, 0, 64)
			mktTmp.AveragePrice, _ = strconv.ParseFloat(averagePrice, 64)
			mktTmp.XSpeedTime = xpeedTime
			mktTmp.TimeStamp = stdTime

			mktDataS = append(mktDataS, mktTmp)
		}

		tmpT, _ = time.Parse(YYYY_MM_DD_hh_mm_ss_offset, time_now)

		//fmt.Println("Unnatual TmpT", tmpT)
		//fmt.Println("Ticks:", mktDataS)

		return mktDataS, tmpT
	}

}

//填充Position数据
func (c *RedisConn) LoadPositionData(subConn RedisSubConn, process ProcessPositionData) error {
	//监听行情并填充行情数据
	if subConn.SubConn == nil {
		return Error{"nil pointer of PubSubConn."}
	}

	var posData PositionData
	var posDataArrayTemp []PositionData

	for {
		switch v := subConn.SubConn.Receive().(type) {
		case redis.Message:
			//处理JSON数据
			data := string(v.Data)
			data = strings.Replace(data, "'", "\"", -1)
			js, err := simplejson.NewJson([]byte(data))

			if err != nil {
				return err
			}

			//填充模型数据
			//number, _ := js.Get("number").String()
			array, _ := js.Get("data").Array()

			posDataArrayTemp = posDataArrayTemp[:0]

			for _, v := range array {
				pos, _ := v.(map[string]interface{})
				accountID := pos["accountID"].(string)
				instrumentID := pos["instrumentID"].(string)
				buySellType := pos["buySellType"].(string)
				totalAvaiAmount := pos["totalAvaiAmount"].(string)
				datePositionProfitLoss := pos["datePositionProfitLoss"].(string)
				dateCloseProfitLoss := pos["dateCloseProfitLoss"].(string)
				isFirst := pos["bIsFirst"].(string)
				isLast := pos["bIsLast"].(string)

				posData.AccountID = accountID
				posData.InstrumentID = instrumentID
				posData.BuySellType = buySellType
				posData.TotalAvaiAmount, _ = strconv.ParseInt(totalAvaiAmount, 0, 64)
				posData.DatePositionProfitLoss, _ = strconv.ParseFloat(datePositionProfitLoss, 64)
				posData.DateCloseProfitLoss, _ = strconv.ParseFloat(dateCloseProfitLoss, 64)
				posData.IsFirst, _ = strconv.ParseInt(isFirst, 0, 64)
				posData.IsLast, _ = strconv.ParseInt(isLast, 0, 64)

				posDataArrayTemp = append(posDataArrayTemp, posData)

			}
			process.ProcessPositionData(posDataArrayTemp)

		case redis.Subscription:
			fmt.Println(v.Channel, v.Count, v.Kind)
		case error:
			fmt.Println(v)
			return v
		}
	}
}

//填充下单响应数据
func (c *RedisConn) LoadPlaceOrderResponseData(subConn RedisSubConn, process ProcessPlaceOrderResponseData) error {
	//监听行情并填充行情数据
	if subConn.SubConn == nil {
		return Error{"nil pointer of PubSubConn."}
	}

	var porData PlaceOrderResponseData
	var porDataArrayTemp []PlaceOrderResponseData

	for {
		switch v := subConn.SubConn.Receive().(type) {
		case redis.Message:

			//处理JSON数据
			data := string(v.Data)
			data = strings.Replace(data, "'", "\"", -1)
			js, err := simplejson.NewJson([]byte(data))

			if err != nil {
				return err
			}

			//填充模型数据
			responseErrorID, _ := js.Get("OnRspInsertOrderErrorID").String()
			localOrderID, _ := js.Get("localOrderID").String()

			porData.ResponseErrorID = responseErrorID
			porData.LocalOrderID, _ = strconv.ParseInt(localOrderID, 0, 64)

			porDataArrayTemp = append(porDataArrayTemp, porData)

			process.ProcessPlaceOrderResponseData(porDataArrayTemp)

		case redis.Subscription:
			fmt.Println(v.Channel, v.Count, v.Kind)
		case error:
			fmt.Println(v)
			return v
		}
	}
}

//填充成交回报数据
func (c *RedisConn) LoadOrderMatchedReturnData(subConn RedisSubConn, process ProcessOrderMatchedReturnData) error {
	//监听行情并填充行情数据
	if subConn.SubConn == nil {
		return Error{"nil pointer of PubSubConn."}
	}

	var omrData OrderMatchedReturnData
	var omrDataArrayTemp []OrderMatchedReturnData

	for {
		switch v := subConn.SubConn.Receive().(type) {
		case redis.Message:
			//处理JSON数据
			data := string(v.Data)
			data = strings.Replace(data, "'", "\"", -1)
			js, err := simplejson.NewJson([]byte(data))

			if err != nil {
				return err
			}

			//填充模型数据
			//instrumentID, _ := js.Get("instrumentID").String() //没有该字段，需清理
			tradePrice, _ := js.Get("TradePrice").String()
			tradeAmount, _ := js.Get("TradeAmount").String()
			tradeTime, _ := js.Get("TradeTime").String()
			localOrderID, _ := js.Get("localOrderID").String()
			matchedTime, _ := js.Get("OnRtnMatchedInfoTime").String()
			//openCloseType, _ := js.Get("open_close_type").String() //只有部分api会返回该字段，需清理

			//omrData.InstrumentID = instrumentID
			omrData.TradePrice, _ = strconv.ParseFloat(tradePrice, 64)
			omrData.TradeAmount, _ = strconv.ParseInt(tradeAmount, 0, 64)
			omrData.TradeTime = tradeTime
			omrData.LocalOrderID, _ = strconv.ParseInt(localOrderID, 0, 64)
			omrData.MatchedTime, _ = strconv.ParseInt(matchedTime, 0, 64)
			//omrData.OpenCloseType = openCloseType

			omrDataArrayTemp = append(omrDataArrayTemp, omrData)

			process.ProcessOrderMatchedReturnData(omrDataArrayTemp)

		case redis.Subscription:
			fmt.Println(v.Channel, v.Count, v.Kind)
		case error:
			fmt.Println(v)
			return v
		}
	}
}

//填充撤单响应数据
func (c *RedisConn) LoadCancelOrderResponseData(subConn RedisSubConn, process ProcessCancelOrderResponseData) error {
	//监听行情并填充行情数据
	if subConn.SubConn == nil {
		return Error{"nil pointer of PubSubConn."}
	}

	var pcoData CancelOrderResponseData
	var pcoDataArrayTemp []CancelOrderResponseData

	for {

		switch v := subConn.SubConn.Receive().(type) {
		case redis.Message:

			//处理JSON数据
			data := string(v.Data)
			data = strings.Replace(data, "'", "\"", -1)

			js, err := simplejson.NewJson([]byte(data))

			if err != nil {
				return err
			}

			//填充模型数据
			responseErrorID, _ := js.Get("OnRspCancelOrderNumber").String()
			localOrderID, _ := js.Get("localOrderID").String()

			pcoData.ResponseErrorID = responseErrorID
			pcoData.LocalOrderID, _ = strconv.ParseInt(localOrderID, 0, 64)

			pcoDataArrayTemp = append(pcoDataArrayTemp, pcoData)

			process.ProcessCancelOrderResponseData(pcoDataArrayTemp)

		case redis.Subscription:
			fmt.Println(v.Channel, v.Count, v.Kind)
		case error:
			fmt.Println(v)
			return v
		}
	}
}

//填充撤单回报数据
func (c *RedisConn) LoadOrderCanceledReturnData(subConn RedisSubConn, process ProcessOrderCanceledReturnData) error {
	//监听行情并填充行情数据
	if subConn.SubConn == nil {
		return Error{"nil pointer of PubSubConn."}
	}

	var pocData OrderCanceledReturnData
	var pocDataArrayTemp []OrderCanceledReturnData

	for {
		switch v := subConn.SubConn.Receive().(type) {
		case redis.Message:

			//处理JSON数据
			data := string(v.Data)
			data = strings.Replace(data, "'", "\"", -1)
			js, err := simplejson.NewJson([]byte(data))

			if err != nil {
				return err
			}

			//填充模型数据
			localOrderID, _ := js.Get("localOrderID").String()
			cancelAmount, _ := js.Get("cancelAmount").String()
			cancelTime, _ := js.Get("canceledTime").String()
			cancelSysTime, _ := js.Get("OnRtnCancelOrderTime").String()

			pocData.LocalOrderID, _ = strconv.ParseInt(localOrderID, 0, 64)
			pocData.CancelAmount, _ = strconv.ParseInt(cancelAmount, 0, 64)
			pocData.CancelTime = cancelTime
			pocData.CancelSysTime, _ = strconv.ParseInt(cancelSysTime, 0, 64)

			pocDataArrayTemp = append(pocDataArrayTemp, pocData)

			process.ProcessOrderCanceledReturnData(pocDataArrayTemp)

		case redis.Subscription:
			fmt.Println(v.Channel, v.Count, v.Kind)
		case error:
			fmt.Println(v)
			return v
		}
	}
}

//填充可用资金数据
func (c *RedisConn) LoadCapitalData(subConn RedisSubConn, process ProcessCapitalData) error {
	//监听行情并填充行情数据
	if subConn.SubConn == nil {
		return Error{"nil pointer of PubSubConn."}
	}

	var cptData CapitalData

	for {
		switch v := subConn.SubConn.Receive().(type) {
		case redis.Message:

			//处理JSON数据
			data := string(v.Data)
			data = strings.Replace(data, "'", "\"", -1)
			js, err := simplejson.NewJson([]byte(data))

			if err != nil {
				return err
			}

			//填充模型数据
			balance, _ := js.Get("Balance").String()
			currentMargin, _ := js.Get("CurrMargin").String()
			available, _ := js.Get("Available").String()

			cptData.Balance, _ = strconv.ParseFloat(balance, 64)
			cptData.CurrentMargin, _ = strconv.ParseFloat(currentMargin, 64)
			cptData.Available, _ = strconv.ParseFloat(available, 64)

			process.ProcessCapitalData(cptData)

		case redis.Subscription:
			fmt.Println(v.Channel, v.Count, v.Kind)
		case error:
			fmt.Println(v)
			return v
		}
	}
}

//下单
//传入一条行情数据，用来获取涨跌停价格处理市价下单情况
func (c *RedisConn) PlaceOrder(chName string, order *OrderData, mktData MarketData) (int64, error) {

	localOrderID, err := c.GetLocalOrderID()
	if err != nil {
		return -1, err
	}

	orderPlaced := *order
	dealPrice := orderPlaced.InsertPrice

	if strings.EqualFold(orderPlaced.OrderPriceType, "marketprice") {

		orderPlaced.OrderPriceType = "limitprice"

		if strings.EqualFold(orderPlaced.BuySellType, "buy") {
			orderPlaced.InsertPrice = mktData.UpperLimitPrice
		}
		if strings.EqualFold(orderPlaced.BuySellType, "sell") {
			orderPlaced.InsertPrice = mktData.LowerLimitPrice
		}
	}

	if strings.EqualFold(orderPlaced.BuySellType, "buy") {
		dealPrice = mktData.AskPrice1
	}
	if strings.EqualFold(orderPlaced.BuySellType, "sell") {
		dealPrice = mktData.BidPrice1
	}

	orderStr := fmt.Sprintf("{\"instrumentCode\":\"%s\",\"instrumentNo\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"orderAmount\":\"%d\",\"orderPriceType\":\"%s\",\"insertPrice\":\"%f\",\"localOrderID\":\"%d\",\"strategy\":\"%s\",\"dealPrice\":\"%f\"}",
		orderPlaced.InstrumentCode,
		orderPlaced.InstrumentNum,
		orderPlaced.BuySellType,
		orderPlaced.OpenCloseType,
		orderPlaced.OrderAmount,
		orderPlaced.OrderPriceType,
		orderPlaced.InsertPrice,
		localOrderID,
		orderPlaced.Strategy,
		dealPrice)

	//下单存储在Redis List中
	if _, err := c.Do("HMSET", "OrderList", localOrderID, orderStr); err != nil {
		return -1, err
	}

	//下单到订单频道
	if _, err := c.Do("PUBLISH", chName, orderStr); err != nil {
		return -1, err
	}
	//c.Flush()
	order.LocalOrderID = localOrderID
	fmt.Println("PUBLISH", chName, orderStr)
	return localOrderID, nil
}

//撤单
func (c *RedisConn) CancelOrder(chName string, calOrder CancelOrderData) error {

	orderStr := fmt.Sprintf("{\"instrumentCode\":\"%s\",\"instrumentNo\":\"%s\",\"localOrderID\":\"%d\"}",
		calOrder.InstrumentCode,
		calOrder.InstrumentNum,
		calOrder.LocalOrderID)

	if _, err := c.Do("PUBLISH", chName, orderStr); err != nil {
		return err
	}
	//c.Flush()
	fmt.Println("PUBLISH", chName, orderStr)
	return nil
}

//LocalOrderID:本地订单号管理
func (c *RedisConn) GetLocalOrderID() (int64, error) {
	v, err := c.Do("INCR", "Global_Local_Order_ID")
	if err != nil {
		return -1, err
	}
	x := v.(int64)
	return x, nil
}

//获取合约交易规则
func GetInstrumentRules(hostNPort string) (map[string]InstrumentRule, error) {
	conn, err := NewConn(hostNPort)
	if conn == nil {
		fmt.Println("GetInstrumentRules():conn = nil.", err)
		return make(map[string]InstrumentRule), err
	}
	defer conn.Close()

	_, err1 := conn.Do("SELECT", "8")
	if err1 != nil {
		return nil, err1
	}

	v, err2 := conn.Do("GET", VAR8_INST_RULES)
	if v == nil {
		return nil, Error{"redis DO GET heyue_rules failed"}
	}
	if err2 != nil {
		return nil, err2
	}

	data := v.([]uint8)
	js, err3 := simplejson.NewJson([]byte(data))

	if err3 != nil {
		return nil, err3
	}

	//填充模型数据
	rules := make(map[string]InstrumentRule)
	tmpMap, _ := js.Map()

	for _, v1 := range tmpMap {

		var tmpRule InstrumentRule

		tmpRule.Code = v1.(map[string]interface{})["code"].(string)
		tmpRule.Name = v1.(map[string]interface{})["name"].(string)

		t1, _ := v1.(map[string]interface{})["unitAmount"].(json.Number).Int64()
		tmpRule.UnitAmount = int(t1)

		months := v1.(map[string]interface{})["months"].([]interface{})
		for _, v2 := range months {
			t2, _ := v2.(json.Number).Int64()
			tmpRule.Months = append(tmpRule.Months, int(t2))
		}

		tmpRule.MarginRatio, _ = v1.(map[string]interface{})["marginRatio"].(json.Number).Float64()
		rules[tmpRule.Code] = tmpRule
	}

	return rules, nil
}

//获取本地持仓
//instID_strategy = instID:strategy
func GetLocalPosition(hostNPort string, instID_strategy string) []LocalPositionData {
	conn, err := redis.Dial("tcp", hostNPort)
	if conn == nil {
		fmt.Println("GetLocalPosition():conn = nil.", err)
		return []LocalPositionData{}
	}
	defer conn.Close()

	//买开和卖开两种
	dataBuy, _ := conn.Do("HGET", instID_strategy, "buy")
	dataSell, _ := conn.Do("HGET", instID_strategy, "sell")

	pos := []LocalPositionData{}
	posBuy := LocalPositionData{}
	posSell := LocalPositionData{}

	if dataBuy != nil {
		j1, err1 := simplejson.NewJson([]byte(dataBuy.([]byte)))
		if err1 != nil {
			fmt.Println(err1)
			return pos
		}

		//填充模型数据
		instID, _ := j1.Get("instrumentID").String()
		bsType, _ := j1.Get("buySellType").String()
		amount, _ := j1.Get("totalAmount").String()
		avgPrice, _ := j1.Get("averagePrice").String()

		posBuy.InstrumentID = instID
		posBuy.BuySellType = bsType
		posBuy.TotalAvaiAmount, _ = strconv.ParseInt(amount, 0, 64)
		posBuy.AveragePrice, _ = strconv.ParseFloat(avgPrice, 64)

		pos = append(pos, posBuy)

	}

	if dataSell != nil {
		j2, err2 := simplejson.NewJson([]byte(dataSell.([]byte)))
		if err2 != nil {
			fmt.Println(err2)
			return pos
		}

		//填充模型数据
		instID, _ := j2.Get("instrumentID").String()
		bsType, _ := j2.Get("buySellType").String()
		amount, _ := j2.Get("totalAmount").String()
		avgPrice, _ := j2.Get("averagePrice").String()

		posSell.InstrumentID = instID
		posSell.BuySellType = bsType
		posSell.TotalAvaiAmount, _ = strconv.ParseInt(amount, 0, 64)
		posSell.AveragePrice, _ = strconv.ParseFloat(avgPrice, 64)

		pos = append(pos, posSell)
	}

	return pos
}

//更新本地持仓
func UpdateLocalPosition(hostNPort string, omr OrderMatchedReturnData) error {
	conn, err := redis.Dial("tcp", hostNPort)
	if conn == nil {
		fmt.Println("UpdateLocalPosition():conn = nil.", err)
		return err
	}
	defer conn.Close()

	//从订单中查询“Buy” OR "Sell"
	order := GetOrderByLocalOrderID(hostNPort, omr.LocalOrderID)

	instrumentID := order.InstrumentCode + order.InstrumentNum
	openCloseType := order.OpenCloseType
	key := instrumentID + ":" + order.Strategy

	localPos := GetLocalPosition(hostNPort, key)

	var posStr string

	if len(localPos) == 0 {
		//没有仓位，直接新增
		if order.OpenCloseType == "open" {
			posStr = fmt.Sprintf("{\"instrumentID\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"totalAmount\":\"%d\",\"averagePrice\":\"%f\",\"strategy\":\"%s\"}",
				instrumentID,
				order.BuySellType,
				openCloseType,
				omr.TradeAmount,
				omr.TradePrice,
				order.Strategy)

			_, err := conn.Do("HMSET", instrumentID+":"+order.Strategy, order.BuySellType, posStr)
			fmt.Println("Updating Local Position [command]:", "HMSET", key, order.BuySellType, posStr)
			if err != nil {
				return err
			}
		}

	} else if len(localPos) == 1 {
		//有一种仓位
		fmt.Println("len = 1")

		if localPos[0].BuySellType == order.BuySellType {
			//开仓
			if order.OpenCloseType == "open" {

				posStr = fmt.Sprintf("{\"instrumentID\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"totalAmount\":\"%d\",\"averagePrice\":\"%f\",\"strategy\":\"%s\"}",
					instrumentID,
					order.BuySellType,
					openCloseType,
					localPos[0].TotalAvaiAmount+omr.TradeAmount,
					Round((localPos[0].AveragePrice*float64(localPos[0].TotalAvaiAmount)+omr.TradePrice*float64(omr.TradeAmount))/float64(localPos[0].TotalAvaiAmount+omr.TradeAmount), 2),
					order.Strategy)

				_, err := conn.Do("HMSET", instrumentID+":"+order.Strategy, order.BuySellType, posStr)
				fmt.Println("Updating Local Position [command]:", "HMSET", key, order.BuySellType, posStr)
				if err != nil {
					return err
				}

			} else if order.OpenCloseType == "closeyesterday" || order.OpenCloseType == "closetoday" {

				fmt.Println("No available position to close.")

			} else {
				fmt.Println("OpenCloseType is not valid, should be open/closeyesterday/closetoday. Getting:", order.OpenCloseType)
				return Error{"OpenCloseType is not valid, should be open/closeyesterday/closetoday."}
			}

		} else {
			//锁仓
			if order.OpenCloseType == "open" {

				posStr = fmt.Sprintf("{\"instrumentID\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"totalAmount\":\"%d\",\"averagePrice\":\"%f\",\"strategy\":\"%s\"}",
					instrumentID,
					order.BuySellType,
					openCloseType,
					omr.TradeAmount,
					omr.TradePrice,
					order.Strategy)

				_, err := conn.Do("HMSET", instrumentID+":"+order.Strategy, order.BuySellType, posStr)
				fmt.Println("Updating Local Position [command]:", "HMSET", key, order.BuySellType, posStr)
				if err != nil {
					return err
				}

			} else if order.OpenCloseType == "closeyesterday" || order.OpenCloseType == "closetoday" {
				//平仓
				posStr = fmt.Sprintf("{\"instrumentID\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"totalAmount\":\"%d\",\"averagePrice\":\"%f\",\"strategy\":\"%s\"}",
					instrumentID,
					localPos[0].BuySellType,
					"open", //剩余仓位，所以依然为open
					localPos[0].TotalAvaiAmount-omr.TradeAmount,
					localPos[0].AveragePrice,
					order.Strategy)

				_, err := conn.Do("HMSET", instrumentID+":"+order.Strategy, localPos[0].BuySellType, posStr)
				fmt.Println("Updating Local Position [command]:", "HMSET", key, localPos[0].BuySellType, posStr)
				if err != nil {
					return err
				}

			} else {
				fmt.Println("OpenCloseType is not valid, should be open/closeyesterday/closetoday. Getting:", order.OpenCloseType)
				return Error{"OpenCloseType is not valid, should be open/closeyesterday/closetoday."}
			}
		}

	} else if len(localPos) == 2 {
		//有两种仓位
		fmt.Println("len = 2")

		var posBuy LocalPositionData
		var posSell LocalPositionData

		for _, v := range localPos {
			if v.BuySellType == "buy" {
				posBuy = v
			} else if v.BuySellType == "sell" {
				posSell = v
			} else {
				fmt.Println("BuySellType in position is not valid, should be buy or sell. Getting:", v)
			}
		}

		//开仓
		if openCloseType == "open" {

			if order.BuySellType == "buy" {

				posStr = fmt.Sprintf("{\"instrumentID\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"totalAmount\":\"%d\",\"averagePrice\":\"%f\",\"strategy\":\"%s\"}",
					instrumentID,
					order.BuySellType,
					openCloseType,
					posBuy.TotalAvaiAmount+omr.TradeAmount,
					Round((posBuy.AveragePrice*float64(posBuy.TotalAvaiAmount)+omr.TradePrice*float64(omr.TradeAmount))/float64(posBuy.TotalAvaiAmount+omr.TradeAmount), 2),
					order.Strategy)

			} else if order.BuySellType == "sell" {

				posStr = fmt.Sprintf("{\"instrumentID\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"totalAmount\":\"%d\",\"averagePrice\":\"%f\",\"strategy\":\"%s\"}",
					instrumentID,
					order.BuySellType,
					openCloseType,
					posSell.TotalAvaiAmount+omr.TradeAmount,
					Round((posSell.AveragePrice*float64(posSell.TotalAvaiAmount)+omr.TradePrice*float64(omr.TradeAmount))/float64(posSell.TotalAvaiAmount+omr.TradeAmount), 2),
					order.Strategy)

			} else {
				fmt.Println("BuySellType in order is not valid, should be buy or sell. Getting:", order)
			}

			_, err := conn.Do("HMSET", instrumentID+":"+order.Strategy, order.BuySellType, posStr)
			fmt.Println("Updating Local Position [command]:", "HMSET", key, order.BuySellType, posStr)
			if err != nil {
				return err
			}

		} else if openCloseType == "closeyesterday" || order.OpenCloseType == "closetoday" {
			//平仓
			if order.BuySellType == "buy" {

				posStr = fmt.Sprintf("{\"instrumentID\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"totalAmount\":\"%d\",\"averagePrice\":\"%f\",\"strategy\":\"%s\"}",
					instrumentID,
					posSell.BuySellType,
					"open",
					posSell.TotalAvaiAmount-omr.TradeAmount,
					posSell.AveragePrice,
					order.Strategy)

				_, err := conn.Do("HMSET", instrumentID+":"+order.Strategy, posSell.BuySellType, posStr)
				fmt.Println("Updating Local Position [command]:", "HMSET", key, posSell.BuySellType, posStr)
				if err != nil {
					return err
				}

			} else if order.BuySellType == "sell" {

				posStr = fmt.Sprintf("{\"instrumentID\":\"%s\",\"buySellType\": \"%s\",\"openCloseType\":\"%s\",\"totalAmount\":\"%d\",\"averagePrice\":\"%f\",\"strategy\":\"%s\"}",
					instrumentID,
					posBuy.BuySellType,
					"open",
					posBuy.TotalAvaiAmount-omr.TradeAmount,
					posBuy.AveragePrice,
					order.Strategy)

				_, err := conn.Do("HMSET", instrumentID+":"+order.Strategy, posBuy.BuySellType, posStr)
				fmt.Println("Updating Local Position [command]:", "HMSET", key, posBuy.BuySellType, posStr)
				if err != nil {
					return err
				}

			} else {
				fmt.Println("BuySellType in order is not valid, should be buy or sell. Getting:", order)
			}

		} else {
			fmt.Println("OpenCloseType is not valid, should be open/closeyesterday/closetoday. Getting:", order.OpenCloseType)
			return Error{"OpenCloseType is not valid, should be open/closeyesterday/closetoday."}
		}

	} else {
		fmt.Println("More than 2 positions of one instrument. Instrument:", instrumentID)
		return Error{"More than 2 positions of one instrument."}
	}

	//清理持仓量为0的记录
	localPos = GetLocalPosition(hostNPort, key)
	for _, v := range localPos {
		if v.TotalAvaiAmount == 0 {
			_, err := conn.Do("HDEL", v.InstrumentID+":"+order.Strategy, v.BuySellType)
			fmt.Println("DELETING:", localPos, "HDEL", v.InstrumentID+":"+order.Strategy, v.BuySellType)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//通过LocalOrderID获取下单信息
func GetOrderByLocalOrderID(hostNPort string, localOrderID int64) OrderData {
	conn, err := redis.Dial("tcp", hostNPort)
	if conn == nil {
		fmt.Println("GetOrderByLocalOrderID():conn = nil.", err)
		return OrderData{}
	}
	defer conn.Close()

	data, err := conn.Do("HGET", "OrderList", localOrderID)

	if err != nil {
		fmt.Println("HGET OrderList Error. msg:", err)
		return OrderData{}
	}

	order := OrderData{}

	if data != nil {
		js, err := simplejson.NewJson([]byte(data.([]byte)))
		if err != nil {
			fmt.Println(err)
		}

		//填充模型数据
		instrumentCode, _ := js.Get("instrumentCode").String()
		instrumentNum, _ := js.Get("instrumentNo").String()
		buySellType, _ := js.Get("buySellType").String()
		openCloseType, _ := js.Get("openCloseType").String()
		orderAmount, _ := js.Get("orderAmount").String()
		orderPriceType, _ := js.Get("orderPriceType").String()
		insertPrice, _ := js.Get("insertPrice").String()
		localOrderID, _ := js.Get("localOrderID").String()
		strategy, _ := js.Get("strategy").String()

		order.InstrumentCode = instrumentCode
		order.InstrumentNum = instrumentNum
		order.BuySellType = buySellType
		order.OpenCloseType = openCloseType
		order.OrderAmount, _ = strconv.ParseInt(orderAmount, 0, 64)
		order.OrderPriceType = orderPriceType
		order.InsertPrice, _ = strconv.ParseFloat(insertPrice, 64)
		order.LocalOrderID, _ = strconv.ParseInt(localOrderID, 0, 64)
		order.Strategy = strategy

		return order
	}
	return OrderData{}
}

//回测主动请求行情全局连接
//单例模式
var backTestPublishConn *RedisConn
var backTestPublishConnOnce sync.Once

func InitBackTestPublishConn(hostNPort string) *RedisConn {

	backTestPublishConnOnce.Do(func() {
		backTestPublishConn, _ = NewConn(hostNPort)
	})
	return backTestPublishConn
}

func GetNextTimeTick() {
	if backTestPublishConn.Conn != nil {
		backTestPublishConn.Do("PUBLISH", CH_TEST_GET_TIME_TICK, "")
	}
}
