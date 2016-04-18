package common

import (
	//"fmt"
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

const (
	HOST = "192.168.99.100:6379"
)

func TestNewConn(t *testing.T) {
	connGood, err1 := NewConn(HOST)
	connBad, err2 := NewConn("0.0.0.0:6379")

	if connGood == nil {
		t.Error("conn should not be nil")
	}
	assert.Equal(t, nil, err1)

	if connBad != nil {
		t.Error("conn should be nil")
	}
	assert.NotEqual(t, nil, err2)
}

func TestGetConnStr(t *testing.T) {
	connGood, _ := NewConn(HOST)
	connStr := connGood.GetConnStr()

	assert.Equal(t, HOST, connStr)
}

func TestSubscribeChan(t *testing.T) {
	connGood1, _ := NewConn(HOST)
	connGood2, _ := NewConn(HOST)
	subConn1, err1 := connGood1.SubscribeChan(CH_MARKET, "SR", "605", true)
	subConn2, err2 := connGood2.SubscribeChan(CH_POSITION, "", "", false)

	assert.NotEqual(t, nil, subConn1)
	assert.Equal(t, nil, err1)
	assert.NotEqual(t, nil, subConn2)
	assert.Equal(t, nil, err2)
}

type proc struct {
	instCode string
	instNum  string

	connMarket  *RedisConn
	connPos     *RedisConn
	connPublish *RedisConn
	connPOR     *RedisConn
	connOMR     *RedisConn
	connPCO     *RedisConn
	connPOC     *RedisConn

	subconnMartet RedisSubConn
	subconnPos    RedisSubConn
	subconnPOR    RedisSubConn
	subconnOMR    RedisSubConn
	subconnPCO    RedisSubConn
	subconnPOC    RedisSubConn

	mkt MarketData        //每笔行情数据
	cdl []CandleStickData //K线数据
	pos PositionData      //持仓数据
	por PlaceOrderResponseData
	omr OrderMatchedReturnData //成交回报数据
	pco CancelOrderResponseData
	poc OrderCanceledReturnData //撤单回报

}

func (p *proc) ProcessMarketData(maketData MarketData) (interface{}, error) {
	//fmt.Println("market:", maketData)
	p.mkt = maketData
	return maketData, nil
}

func (p *proc) ProcessCandleStickData(candleData []CandleStickData) (interface{}, error) {

	p.cdl = candleData
	//fmt.Println(candleData[len(candleData)-1])
	return nil, nil
}

func (p *proc) ProcessPositionData(posData []PositionData) (interface{}, error) {
	//fmt.Println("position:", posData)
	p.pos = posData[len(posData)-1]
	return posData, nil
}

func (p *proc) ProcessPlaceOrderResponseData(porData []PlaceOrderResponseData) (interface{}, error) {
	//fmt.Println("placedOrderResponse:", porData)
	p.por = porData[len(porData)-1]
	return porData, nil
}

func (p *proc) ProcessOrderMatchedReturnData(omrData []OrderMatchedReturnData) (interface{}, error) {
	//fmt.Println("matchedOrderReturn:", omrData)
	p.omr = omrData[len(omrData)-1]
	return omrData, nil
}

func (p *proc) ProcessCancelOrderResponseData(pcoData []CancelOrderResponseData) (interface{}, error) {
	//fmt.Println("canceledOrderResponse:", pcoData)
	p.pco = pcoData[len(pcoData)-1]
	return pcoData, nil
}

func (p *proc) ProcessOrderCanceledReturnData(pocData []OrderCanceledReturnData) (interface{}, error) {
	//fmt.Println("canceledOrderReturn:", pocData)
	p.poc = pocData[len(pocData)-1]
	return pocData, nil
}

var p proc

func Init() {
	var proc1 proc

	proc1.instCode = "y"
	proc1.instNum = "1605"

	proc1.connMarket, _ = NewConn(HOST)
	proc1.connPos, _ = NewConn(HOST)
	proc1.connPublish, _ = NewConn(HOST)
	proc1.connPOR, _ = NewConn(HOST)
	proc1.connOMR, _ = NewConn(HOST)
	proc1.connPCO, _ = NewConn(HOST)
	proc1.connPOC, _ = NewConn(HOST)

	proc1.subconnMartet, _ = proc1.connMarket.SubscribeChan(CH_MARKET, proc1.instCode, proc1.instNum, true)
	proc1.subconnPos, _ = proc1.connPos.SubscribeChan(CH_POSITION, "", "", false)
	proc1.subconnPOR, _ = proc1.connPOR.SubscribeChan(CH_INSERT_ORDER_RSP, "", "", false)
	proc1.subconnOMR, _ = proc1.connOMR.SubscribeChan(CH_INSERT_ORDER_MATCHED_RTN, "", "", false)
	proc1.subconnPCO, _ = proc1.connPCO.SubscribeChan(CH_CANCEL_ORDER_RSP, "", "", false)
	proc1.subconnPOC, _ = proc1.connPOC.SubscribeChan(CH_CANCEL_ORDER_RTN, "", "", false)

	p = proc1

}

func TestLoadMarketData(t *testing.T) {

	go p.connMarket.LoadMarketData(p.subconnMartet, &p, 0)
	time.Sleep(time.Second * 2)

	assert.Equal(t, "y1605", p.mkt.InstrumentID)
}

func TestLoadPositionData(t *testing.T) {

	go p.connPos.LoadPositionData(p.subconnPos, &p)
	p.connPublish.Do("PUBLISH", CH_POSITION_PUBLISH, "")
	time.Sleep(time.Second * 1)

	assert.Equal(t, "y1605", p.pos.InstrumentID)
}

//同样测试了PlaceOrder函数
func TestLoadPlaceOrderResponseData(t *testing.T) {

	go p.connPOR.LoadPlaceOrderResponseData(p.subconnPOR, &p)

	p.connPublish.PlaceOrder(CH_INSERT_ORDER, &OrderData{"y", "1605", "buy", "open", 1, "marketprice", 12000, 1, "test"}, p.mkt)
	time.Sleep(time.Second * 1)

	assert.Equal(t, "0", p.por.ResponseErrorID)
}

func TestLoadOrderMatchedReturnData(t *testing.T) {

	go p.connOMR.LoadOrderMatchedReturnData(p.subconnOMR, &p)
	time.Sleep(time.Second * 1)

	p.connPublish.PlaceOrder(CH_INSERT_ORDER, &OrderData{"CF", "605", "buy", "open", 1, "marketprice", 12000, 1, "test"}, p.mkt)
	time.Sleep(time.Second * 3)

	assert.Equal(t, "CF605", p.omr.InstrumentID)
	assert.Equal(t, int64(1), p.omr.TradeAmount)
}

//同样测试了CancelOrder函数
func TestLoadCancelOrderResponseData(t *testing.T) {

	go p.connPCO.LoadCancelOrderResponseData(p.subconnPCO, &p)

	localOrderID, _ := p.connPublish.PlaceOrder(CH_INSERT_ORDER, &OrderData{"y", "1605", "buy", "open", 1, "limitprice", 11600, 1, "test"}, p.mkt)
	time.Sleep(time.Second * 1)
	//撤单
	p.connPublish.CancelOrder(CH_CANCEL_ORDER, CancelOrderData{"y", "1605", localOrderID})
	time.Sleep(time.Second * 1)

	assert.Equal(t, "0", p.pco.ResponseErrorID)
}

func TestLoadOrderCanceledReturnData(t *testing.T) {

	go p.connPOC.LoadOrderCanceledReturnData(p.subconnPOC, &p)

	localOrderID, _ := p.connPublish.PlaceOrder(CH_INSERT_ORDER, &OrderData{"y", "1605", "buy", "open", 1, "limitprice", 11600, 1, "test"}, p.mkt)
	time.Sleep(time.Second * 1)
	//撤单
	p.connPublish.CancelOrder(CH_CANCEL_ORDER, CancelOrderData{"y", "1605", localOrderID})
	time.Sleep(time.Second * 1)

	assert.Equal(t, localOrderID, p.poc.LocalOrderID)
	assert.Equal(t, int64(1), p.poc.CancelAmount)
}

func TestLocalOrderID(t *testing.T) {
	connGood1, _ := NewConn(HOST)
	connGood2, _ := NewConn(HOST)

	id1, _ := connGood1.GetLocalOrderID()
	id2, _ := connGood2.GetLocalOrderID()
	assert.Equal(t, id1+1, id2)
}
