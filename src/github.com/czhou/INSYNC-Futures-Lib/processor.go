package common

import (
	"log"
	"strings"
	"time"
)

type Proc struct {
	InstCode string //合约代码
	InstNum  string //合约期数
	StgName  string //策略实例名
	ConnStr  string //连接字符串
	IsNatual bool   //是否是接收自然行情（非回测行情）

	ConnMarket  *RedisConn //行情连接
	ConnPOS     *RedisConn //持仓连接
	ConnPublish *RedisConn //发送Redis命令连接
	ConnOMR     *RedisConn //成交回报
	ConnPOC     *RedisConn //撤单回报

	SubconnMartet RedisSubConn //行情订阅频道
	SubconnPOS    RedisSubConn //持仓订阅频道
	SubconnOMR    RedisSubConn //成交回报频道
	SubconnPOC    RedisSubConn //撤单回报频道

	Mkt           []MarketData                        //每笔行情数据
	LastTimeStamp time.Time                           //最新一笔行情的时间戳
	Cdl           map[time.Duration][]CandleStickData //K线数据
	KParams       []KParams                           //K线参数

	LocalOrderIDs []OrderIDStatus     //本地下单状态
	LocalPos      []LocalPositionData //本地持仓数据

	Rule       InstrumentRule //交易规则
	CloseTimes []time.Time    //闭市时间

	InProcessMarketData              func(this *Proc, marketData MarketData) (interface{}, error)             //行情触发回调函数
	InProcessCandleStickData         func(this *Proc, candleData []CandleStickData) (interface{}, error)      //K线触发回调函数
	InProcessPositionData            func(this *Proc, posData []PositionData) (interface{}, error)            //服务器持仓触发回调函数
	InProcessOrderMatchedReturnData  func(this *Proc, omrData []OrderMatchedReturnData) (interface{}, error)  //成交回报触发回调函数
	InProcessOrderCanceledReturnData func(this *Proc, pocData []OrderCanceledReturnData) (interface{}, error) //撤单回报触发回调函数

}

type OrderIDStatus struct {
	localOrderID int64
	amount       int64
	dealedAmount int64
	status       int
}

//行情回调函数
func (p *Proc) ProcessMarketData(marketData MarketData) (interface{}, error) {

	if p.InstCode+p.InstNum == marketData.InstrumentID {
		p.Mkt = append(p.Mkt, marketData)

		tmpT, _ := time.Parse(YYYY_MM_DD_hh_mm_ss_MST, marketData.TimeStamp+" CST")
		p.LastTimeStamp = tmpT

		if p.InProcessMarketData != nil {
			p.InProcessMarketData(p, marketData)
		}
	}

	return nil, nil
}

//K线回调函数
func (p *Proc) ProcessCandleStickData(candleData []CandleStickData) (interface{}, error) {

	if p.InstCode+p.InstNum == candleData[len(candleData)-2].InstrumentID {
		kDuration := candleData[len(candleData)-2].KDuration
		p.Cdl[kDuration] = append(p.Cdl[kDuration], candleData[len(candleData)-2])

		if p.InProcessCandleStickData != nil {
			p.InProcessCandleStickData(p, candleData)
		}
	}

	return nil, nil
}

//持仓回调函数
func (p *Proc) ProcessPositionData(posData []PositionData) (interface{}, error) {

	if p.InProcessPositionData != nil {
		p.InProcessPositionData(p, posData)
	}

	return nil, nil
}

//成交回报回调函数
func (p *Proc) ProcessOrderMatchedReturnData(omrData []OrderMatchedReturnData) (interface{}, error) {

	for i, v := range p.LocalOrderIDs {
		if v.status == ORDER_STATUS_DONE {
			continue
		} else {
			if v.localOrderID == omrData[len(omrData)-1].LocalOrderID {
				//更新订单状态
				p.LocalOrderIDs[i].dealedAmount += omrData[len(omrData)-1].TradeAmount
				if p.LocalOrderIDs[i].dealedAmount < v.amount {
					p.LocalOrderIDs[i].status = ORDER_STATUS_PART_DONE
				} else {
					p.LocalOrderIDs[i].status = ORDER_STATUS_DONE
				}

				err := UpdateLocalPosition(p.ConnMarket.connStr, omrData[len(omrData)-1])
				if err != nil {
					log.Println(err)
				}

				//更新本地持仓
				p.LocalPos = GetLocalPosition(p.ConnMarket.connStr, p.InstCode+p.InstNum+":"+p.StgName)

				if p.InProcessOrderMatchedReturnData != nil {
					p.InProcessOrderMatchedReturnData(p, omrData)
				}
			}
		}
	}

	return nil, nil
}

//撤单回报回调函数
func (p *Proc) ProcessOrderCanceledReturnData(pocData []OrderCanceledReturnData) (interface{}, error) {

	if p.InProcessOrderCanceledReturnData != nil {
		p.InProcessOrderCanceledReturnData(p, pocData)
	}

	return nil, nil
}

//下单函数
func (p *Proc) PlaceOrder(buySellType, openCloseType string, price float64, amount int64, priceType string) (localOrderID int64, err error) {

	if !strings.EqualFold(buySellType, "buy") && !strings.EqualFold(buySellType, "sell") {
		log.Println("Place Order Failed, invalid buySellType", buySellType)
		return -1, Error{"Place Order Failed, invalid buySellType"}
	}

	if !strings.EqualFold(openCloseType, "open") && !strings.EqualFold(openCloseType, "close") && !strings.EqualFold(openCloseType, "closetoday") {
		log.Println("Place Order Failed, invalid openCloseType", openCloseType)
		return -1, Error{"Place Order Failed, invalid openCloseType"}
	}

	if !strings.EqualFold(priceType, "limitprice") && !strings.EqualFold(priceType, "marketprice") {
		log.Println("Place Order Failed, invalid priceType", priceType)
		return -1, Error{"Place Order Failed, invalid priceType"}
	}

	//还没有收到一条行情
	if len(p.Mkt) == 0 {
		log.Println("Place Order Failed, no market data recieved yet, instID:", p.InstCode+p.InstNum)
		return -1, Error{"Place Order Failed, no market data recieved yet."}
	}

	//下单
	order := OrderData{p.InstCode, p.InstNum, buySellType, openCloseType, amount, priceType, price, 0, p.StgName}
	localOrderID, err = p.ConnPublish.PlaceOrder(CH_INSERT_ORDER, &order, p.Mkt[len(p.Mkt)-1])

	//记录单号
	if err == nil {
		p.LocalOrderIDs = append(p.LocalOrderIDs, OrderIDStatus{localOrderID, amount, 0, ORDER_STATUS_PLACED})
	}

	return localOrderID, err
}

//撤单函数
func (p *Proc) CancelOrder(localOrderID int64) error {

	//检查单号
	isExist := false

	for i, v := range p.LocalOrderIDs {
		if v.localOrderID != localOrderID {
			continue
		} else if v.status != ORDER_STATUS_PLACED || v.status != ORDER_STATUS_PART_DONE {
			isExist = true
			log.Println("Cancel Order Failed, order is already done. LocalOrderID:", localOrderID)
			return Error{"Cancel Order Failed, order is already done."}
		} else {
			//撤单
			order := CancelOrderData{p.InstCode, p.InstNum, localOrderID}
			err := p.ConnPublish.CancelOrder(CH_CANCEL_ORDER, order)

			if err == nil {
				if v.dealedAmount == 0 {
					//完全撤单
					p.LocalOrderIDs[i].status = ORDER_STATUS_CANCELED

				} else {
					//部分撤单
					p.LocalOrderIDs[i].status = ORDER_STATUS_PART_CANCELED
				}
			}

		}
	}

	//单号不存在
	if !isExist {
		log.Println("Cancel Order Failed, order does not exist. LocalOrderID:", localOrderID)
		return Error{"Cancel Order Failed, order does not exist."}
	}

	return nil
}

func NewProc(hostNPort, instCode, instNum, strategy string, kParams []KParams, isNatual bool) *Proc {
	var proc1 Proc

	proc1.InstCode = instCode
	proc1.InstNum = instNum
	proc1.StgName = strategy
	proc1.ConnStr = hostNPort
	proc1.IsNatual = isNatual

	//初始化系统时间
	sysT := GetSysTime()

	//实例化K线map
	proc1.Cdl = make(map[time.Duration][]CandleStickData)
	proc1.KParams = kParams

	proc1.ConnMarket, _ = NewConn(hostNPort)
	proc1.ConnPOS, _ = NewConn(hostNPort)
	proc1.ConnPublish, _ = NewConn(hostNPort)
	proc1.ConnOMR, _ = NewConn(hostNPort)
	proc1.ConnPOC, _ = NewConn(hostNPort)

	if isNatual {
		proc1.SubconnMartet, _ = proc1.ConnMarket.SubscribeChan(CH_MARKET, proc1.InstCode, proc1.InstNum, true)
	} else {
		proc1.SubconnMartet, _ = proc1.ConnMarket.SubscribeChan(CH_TEST_MARKET, proc1.InstCode, proc1.InstNum, false)
	}
	proc1.SubconnMartet.IsUnnatual = !isNatual

	proc1.SubconnPOS, _ = proc1.ConnPOS.SubscribeChan(CH_POSITION, "", "", false)
	proc1.SubconnOMR, _ = proc1.ConnOMR.SubscribeChan(CH_INSERT_ORDER_MATCHED_RTN, "", "", false)
	proc1.SubconnPOC, _ = proc1.ConnPOC.SubscribeChan(CH_CANCEL_ORDER_RTN, "", "", false)

	go proc1.ConnMarket.LoadMarketData(proc1.SubconnMartet, &proc1, kParams)
	go proc1.ConnPOS.LoadPositionData(proc1.SubconnPOS, &proc1)
	go proc1.ConnOMR.LoadOrderMatchedReturnData(proc1.SubconnOMR, &proc1)
	go proc1.ConnPOC.LoadOrderCanceledReturnData(proc1.SubconnPOC, &proc1)

	//初始化合约交易规则
	rules, e := GetInstrumentRules(hostNPort)
	if e != nil {
		log.Println(e)
	}
	proc1.Rule = rules[proc1.InstCode]

	//初始化合约的闭市时间
	tradingTimes, _ := GetTradingTimeSlots(proc1.InstCode, proc1.IsNatual)
	if len(tradingTimes) == 0 {
		log.Println("Instrument trading time is not initialized.")
	} else {
		timeNow := sysT.Now(proc1.IsNatual).In(time.FixedZone("CST", 28800))
		proc1.CloseTimes = append(proc1.CloseTimes, tradingTimes[len(tradingTimes)-1])

		if tradingTimes[len(tradingTimes)-1].Truncate(time.Hour*24).Sub(tradingTimes[0].Truncate(time.Hour*24)) > 0 {
			dayStart := timeNow.Truncate(time.Hour * 24).Add(-time.Hour * 8)
			nightTradeEndP, _ := time.Parse("15:04:05", NIGHT_C1_TRADE_END_P1)
			//获取的夜盘收盘时间是前一天，所以+1天
			proc1.CloseTimes = append(proc1.CloseTimes, dayStart.Add(time.Duration(nightTradeEndP.Hour()-24)*time.Hour).Add(time.Duration(nightTradeEndP.Minute())*time.Minute).Add(time.Hour*24))
		}

		if len(proc1.CloseTimes) == 1 && timeNow.Hour() > 15 {
			proc1.CloseTimes[0] = proc1.CloseTimes[0].Add(-time.Hour * 24)
		}
	}

	return &proc1
}

//资金结构体
type Capital struct {
	ConnCPT    *RedisConn   //可用资金连接
	SubconnCPT RedisSubConn //可用资金频道

	Cpt                CapitalData
	InitAvailable      float64 //初始可用资金量
	IsInitAvailableSet bool    //初始可用资金是否已被初始化
}

func (p *Capital) ProcessCapitalData(cptData CapitalData) (interface{}, error) {
	//初始化可用资金量
	if !p.IsInitAvailableSet {
		p.InitAvailable = cptData.Available
		p.IsInitAvailableSet = true
	}

	p.Cpt = cptData
	log.Println("capital:", cptData)
	return nil, nil
}

func NewCapital(hostNPort string) *Capital {
	var capital Capital

	capital.ConnCPT, _ = NewConn(hostNPort)
	capital.SubconnCPT, _ = capital.ConnCPT.SubscribeChan(CH_CAPITAL, "", "", false)

	go capital.ConnCPT.LoadCapitalData(capital.SubconnCPT, &capital)

	return &capital
}
