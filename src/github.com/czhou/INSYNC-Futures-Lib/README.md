# INSYNC-Futures-StrategyCommonLib


## USAGE:
package main

import (
	"common"
	"fmt"
)

func main() {
	
	//创建阻塞信道，保持main函数不退出
	ch := make(chan bool)

	//所有订阅行情的连接必须是独立的连接
	//行情数据连接
	connMarket, _ := common.NewConn("192.168.59.103:6379")
	//持仓数据连接
	connPos, _ := common.NewConn("192.168.59.103:6379")
	//公用Publish连接，用来向Redis发送Publish命令
	connPublish, _ := common.NewConn("192.168.59.103:6379")
	//下单响应连接
	connPOR, _ := common.NewConn("192.168.59.103:6379")
	//成交回报连接
	connOMR, _ := common.NewConn("192.168.59.103:6379")
	//撤单响应连接
	connPCO, _ := common.NewConn("192.168.59.103:6379")
	//撤单回报连接
	connPOC, _ := common.NewConn("192.168.59.103:6379")
	//资金数据连接
	connCPT, _ := common.NewConn("192.168.99.100:6379")

	//订阅行情频道
	subconnMartet, _ := connMarket.SubscribeChan(common.CH_MARKET, "SR", "601", true)
	//订阅持仓频道
	subconnPos, _ := connPos.SubscribeChan(common.CH_POSITION, "", "", false)
	//订阅下单响应频道
	subconnPOR, _ := connPOR.SubscribeChan(common.CH_INSERT_ORDER_RSP, "", "", false)
	//订阅下单回报频道
	subconnOMR, _ := connOMR.SubscribeChan(common.CH_INSERT_ORDER_MATCHED_RTN, "", "", false)
	//订阅撤单响应频道
	subconnPCO, _ := connPCO.SubscribeChan(common.CH_CANCEL_ORDER_RSP, "", "", false)
	//订阅撤单回报频道
	subconnPOC, _ := connPOC.SubscribeChan(common.CH_CANCEL_ORDER_RTN, "", "", false)
	//订阅资金数据频道
	subconnCPT, _ := connCPT.SubscribeChan(common.CH_CAPITAL, "", "", false)

	//声明一个mar结构体变量，mar结构体参见main()函数以下
	var mkt mar
	
	//初始化需要获取的数据种类
	kprams1 := []common.KParams{{common.K_NONE, 0}}	//只获取分时行情数据
	kprams3 := []common.KParams{{common.K_NONE, 0}, {common.K_SIMPLE, time.Second * 5}, {common.K_COMPLEX, time.Second * 10}}	//获取分时、简单5秒钟K线，复杂10秒钟K线

	//启动不同的协程来处理订阅频道的数据
	//逐笔接收处理行情数据，回调函数方式
	go connMarket.LoadMarketData(subconnMartet, &mkt, kprams1)
	//接收处理指定时间间隔的K线数据，回调函数方式
	go connMarket.LoadMarketData(subconnMartet, &mkt, kprams3)
	
	//非回调函数方式获取K线数据，mar1结构体不同于上面的mar结构体
	var mkt mar1
	go connMarket.LoadMarketData(subconnMartet, &mkt, kprams1)

	//接受处理持仓数据
	go connPos.LoadPositionData(subconnPos, &mkt)
	//接收处理下单响应数据
	go connPOR.LoadPlaceOrderResponseData(subconnPOR, &mkt)
	//接收处理成交回报数据
	go connOMR.LoadOrderMatchedReturnData(subconnOMR, &mkt)
	//接收处理撤单响应数据
	go connPCO.LoadCancelOrderResponseData(subconnPCO, &mkt)
	//接收处理撤单回报数据
	go connPOC.LoadOrderCanceledReturnData(subconnPOC, &mkt)
	//接收资金数据
	go connCPT.LoadCapitalData(subconnCPT, &mkt)


	//查询持仓
	connPublish.Do("PUBLISH", common.CH_POSITION_PUBLISH, "")
	//查询资金
	connPublish.Do("PUBLISH", common.CH_CAPITAL_PUBLISH, "")
	
	//下单
	//localOrderID可以填任何数，下单成功后结构体的localOrderID会被设置，并且下单函数会返回localOrderID
	odr := common.OrderData{"CF", "601", "buy", "open", 1, "limitprice", 12200, 0, "test"}
	localOrderID, _ := connPublish.PlaceOrder(common.CH_INSERT_ORDER, &odr)

	//撤单
	//利用下单函数返回的localOrderID来撤单
	connPublish.CancelOrder(common.CH_CANCEL_ORDER, common.CancelOrderData{"CF", "601", localOrderID})

	//阻塞main
	<-ch
}

//结构体mar，可以根据不同的需求自行定义  
//它实现了处理订阅频道返回的接口，接口函数自行实现  
type mar struct {
}

//eg:  
//处理行情数据  
func (p *mar) ProcessMarketData(maketData common.MarketData) (interface{}, error) {
	fmt.Println("market:", maketData)
	return nil, nil
}

//处理K线数据
//回调函数方式，在下一条K线的第一条行情到来时调用，所以处理的K线数据是candleData[len(candleData)-2]
func (p *mar) ProcessCandleStickData(candleData []common.CandleStickData) (interface{}, error) {
	
	//处理5秒钟k线
	if candleData[len(candleData)-1].KDuration == time.Second*5 {
		fmt.Println("Duration[5s]:", candleData[len(candleData)-2])
	}

	//处理10秒钟k线
	if candleData[len(candleData)-1].KDuration == time.Second*10 {
		fmt.Println("Duration[10s]:", candleData[len(candleData)-2])
	}
	
	return nil, nil
}
//非回调函数方式
func (p *mar) ProcessMarketData(maketData common.MarketData) (interface{}, error) {

	fmt.Println(maketData)
	//获取一天的K线时间切片
	//传入品种代码和K线周期
	kTimeSlots, err := common.GetIntervalTimeSlots("ag", time.Second*5)

	if err != nil {
		return nil, err
	}

	p.kData = common.GetCandleStickData(kTimeSlots, p.kData, maketData)

	for i, v := range p.kData {
		fmt.Printf("[CandleData][%d]:%v\n", i, v)
	}

	return nil, nil
}


//处理持仓数据  
func (p *mar) ProcessPositionData(posData []common.PositionData) (interface{}, error) {
	fmt.Println("position:", posData)
	return nil, nil
}

//处理下单响应  
func (p *mar) ProcessPlaceOrderResponseData(porData []common.PlaceOrderResponseData) (interface{}, error) {
	fmt.Println("placedOrderResponse:", porData)
	return nil, nil
}

//处理成交回报  
func (p *mar) ProcessOrderMatchedReturnData(omrData []common.OrderMatchedReturnData) (interface{}, error) {
	fmt.Println("matchedOrderReturn:", omrData)
	return nil, nil
}

//处理撤单响应  
func (p *mar) ProcessCancelOrderResponseData(pcoData []common.CancelOrderResponseData) (interface{}, error) {
	fmt.Println("canceledOrderResponse:", pcoData)
	return nil, nil
}

//处理撤单回报  
func (p *mar) ProcessOrderCanceledReturnData(pocData []common.OrderCanceledReturnData) (interface{}, error) {
	fmt.Println("canceledOrderReturn:", pocData)
	return nil, nil
}

//处理资金数据
func (p *mar) ProcessCapitalData(cptData common.CapitalData) (interface{}, error) {
	fmt.Println("capital:", cptData)
	return nil, nil
}



//非回调函数方式获取K线数据
type mar1 struct {
	kData []common.CandleStickData //简单K线数据
	kComplexData []common.CandleStickData //复杂K线数据
}

func (p *mar1) ProcessMarketData(maketData common.MarketData) (interface{}, error) {

	//获取一天的K线时间切片
	//传入品种代码和K线周期
	kTimeSlots, err := common.GetIntervalTimeSlots("ag", time.Second*5)

	if err != nil {
		return nil, err
	}

	p.kData = common.GetCandleStickData(kTimeSlots, p.kData, maketData)
	p.kComplexData = common.GetComplexCandleStickData(kTimeSlots, p.kComplexData, maketData)

	for i, v := range p.kData {
		fmt.Printf("[CandleData][%d]:%v\n", i, v)
	}

	return nil, nil
}

func (p *mar1) ProcessCandleStickData(candleData []common.CandleStickData) (interface{}, error) {

	fmt.Println(candleData[len(candleData)-2])
	return nil, nil
}



## 频道信息
DB8:futures:MarketData		//行情频道  
futures:pos2				//持仓数据返回  
DB8:futures:ReqQryPosition	//查询持仓 每次返回一条  
DB8:futures:ReqQryPosition2	//查询持仓 一次返回持仓信息
DB8:futures:ReqQryPosition3	//查询持仓 一次返回本地策略持仓

DB8:futures:ReqInsertOrder	//下单  
DB8:futures:ReqCancelOrder	//撤单

futures:OnRspInsertOrder	//下单响应  
futures:OnRtnOrder 			//委托回报  
futures:OnRtnMatchedInfo	//成交回报  
futures:OnRspCancelOrder	//撤单响应  
futures:OnRtnCancelOrder	//撤单回报