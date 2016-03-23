package main

import (
	"flag"
	"github.com/czhou/INSYNC-Futures-Lib"
	"log"
	"time"
)

func main() {
	//hostNPort := "192.168.99.100:6379"

	hostNPort := flag.String("h", "hostNPort", "String contains host and port. Format 192.168.1.1:8888")
	strategyInst := flag.String("s", "strategyInst", "strategy instance name")
	flag.Parse()

	log.Println(*strategyInst)
	kprams3 := []common.KParams{{common.K_NONE, time.Duration(0)}, {common.K_COMPLEX, time.Minute * 10}} //获取分时、简单5秒钟K线，复杂10秒钟K线

	proc1 := common.NewProc(*hostNPort, "rb", "1610", "YYYY", kprams3, true)
	proc1.InProcessMarketData = processMarketData
	proc1.InProcessCandleStickData = processCandleStickData

	//common.InitBackTestPublishConn(hostNPort)

	//	i := 0
	//	t1 := time.Now()
	//	for i = 0; i < 10000; i++ {
	//		//time.Sleep(time.Second * 1)
	//		common.GetNextTimeTick()
	//	}
	//	t2 := time.Now()
	//	log.Println("i:=", i, "time:", t2.Sub(t1))

	//程序退出
	exitTime := time.NewTimer(time.Hour * 100)
	log.Println("program exited at:", <-exitTime.C)

}

type mar struct {
	kData []common.CandleStickData //K线数据
	cpt   common.CapitalData
}

func (p *mar) ProcessMarketData(maketData common.MarketData) (interface{}, error) {

	log.Println("[MarketData]:", maketData)
	return nil, nil
}

func (p *mar) ProcessCandleStickData(candleData []common.CandleStickData) (interface{}, error) {

	log.Println("Kdata:", candleData[len(candleData)-2])

	return nil, nil
}

//函数变量
var processMarketData = func(this *common.Proc, marketData common.MarketData) (interface{}, error) {
	log.Println("Market:", marketData)
	return nil, nil
}

var processCandleStickData = func(this *common.Proc, candleData []common.CandleStickData) (interface{}, error) {
	log.Println("Candle:", candleData[len(candleData)-1])
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
