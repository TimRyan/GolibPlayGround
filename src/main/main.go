package main

import (
	"flag"
	"github.com/czhou/INSYNC-Futures-Lib"
	"log"
	"time"
)

func main() {

	hostNPort := flag.String("h", "hostNPort", "String contains host and port. Format 192.168.1.1:8888")
	strategyInst := flag.String("s", "strategyInst", "strategy instance name")
	flag.Parse()

	kprams := []common.KParams{{common.K_NONE, 0}, {common.K_SIMPLE, time.Second * 5}, {common.K_COMPLEX, time.Second * 10}} //获取分时、简单5秒钟K线，复杂10秒钟K线
	proc := common.NewProc(*hostNPort, "l", "1605", *strategyInst, kprams, true)
	proc.InProcessMarketData = processMarketData
	proc.InProcessCandleStickData = processCandleStickData

	//	kprams3 := []common.KParams{{common.K_NONE, time.Duration(0)}} //获取分时、简单5秒钟K线，复杂10秒钟K线

	//	proc1 := common.NewProc("192.168.99.100:6379", "l", "1605", "YYYY", kprams3, true)
	//	proc1.InProcessMarketData = processMarketData
	//	proc1.InProcessOrderMatchedReturnData = processOrderMatchedReturnData

	//	//测试下单和平仓

	//	//买开
	//	time.Sleep(time.Second * 2)
	//	proc1.PlaceOrder("buy", "open", 0, 5, "marketprice")
	//	time.Sleep(time.Second * 5)
	//	log.Println("LocalIDs--->:", proc1.LocalOrderIDs, "LocalPos:", proc1.LocalPos)
	//	proc1.PlaceOrder("buy", "open", 0, 5, "marketprice")
	//	time.Sleep(time.Second * 5)
	//	log.Println("LocalIDs--->:", proc1.LocalOrderIDs, "LocalPos:", proc1.LocalPos)
	//	//卖开
	//	proc1.PlaceOrder("sell", "open", 0, 5, "marketprice")
	//	time.Sleep(time.Second * 5)
	//	log.Println("LocalIDs--->:", proc1.LocalOrderIDs, "LocalPos:", proc1.LocalPos)
	//	proc1.PlaceOrder("sell", "open", 0, 5, "marketprice")
	//	time.Sleep(time.Second * 5)
	//	log.Println("LocalIDs--->:", proc1.LocalOrderIDs, "LocalPos:", proc1.LocalPos)

	//	//买平
	//	proc1.PlaceOrder("buy", "close", 0, 5, "marketprice")
	//	time.Sleep(time.Second * 5)
	//	log.Println("LocalIDs--->:", proc1.LocalOrderIDs, "LocalPos:", proc1.LocalPos)
	//	proc1.PlaceOrder("buy", "close", 0, 5, "marketprice")
	//	time.Sleep(time.Second * 5)
	//	log.Println("LocalIDs--->:", proc1.LocalOrderIDs, "LocalPos:", proc1.LocalPos)
	//	//卖平
	//	proc1.PlaceOrder("sell", "close", 0, 5, "marketprice")
	//	time.Sleep(time.Second * 5)
	//	log.Println("LocalIDs--->:", proc1.LocalOrderIDs, "LocalPos:", proc1.LocalPos)
	//	proc1.PlaceOrder("sell", "close", 0, 5, "marketprice")
	//	time.Sleep(time.Second * 5)
	//	log.Println("LocalIDs--->:", proc1.LocalOrderIDs, "LocalPos:", proc1.LocalPos)

	//程序退出
	exitTime := time.NewTimer(time.Hour * 100)
	log.Println("program exited at:", <-exitTime.C)

}

//函数变量
var processMarketData = func(this *common.Proc, marketData common.MarketData) (interface{}, error) {
	log.Println("Market:", marketData)
	return nil, nil
}

var processCandleStickData = func(this *common.Proc, candleData []common.CandleStickData) (interface{}, error) {
	log.Println("Candle:", candleData[len(candleData)-2])
	return nil, nil
}

var processPositionData = func(this *common.Proc, posData []common.PositionData) (interface{}, error) {
	//fmt.Println("posData:", posData)
	return nil, nil
}

var processOrderMatchedReturnData = func(this *common.Proc, omrData []common.OrderMatchedReturnData) (interface{}, error) {
	//log.Println("OMR:", omrData)
	return nil, nil
}

var processOrderCanceledReturnData = func(this *common.Proc, pocData []common.OrderCanceledReturnData) (interface{}, error) {
	//fmt.Println("POC:", pocData)
	return nil, nil
}
