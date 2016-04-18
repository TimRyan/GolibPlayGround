package common

import (
	"time"
)

//行情数据
type MarketData struct {
	TradingDate        string  //交易日
	InstrumentID       string  //合约代码
	ExchangeID         string  //交易所代码
	ExchangeInstID     string  //合约在交易所的代码
	LastPrice          float64 //最新价
	PreSettlementPrice float64 //上次结算价
	PreClosePrice      float64 //昨收盘
	PreOpenInterest    int64   //昨持仓量
	OpenPrice          float64 //今开盘
	HighestPrice       float64 //最高价
	LowestPrice        float64 //最低价
	Volume             int64   //成交数量
	Turnover           float64 //成交金额
	OpenInterest       int64   //持仓量
	ClosePrice         float64 //今收盘
	SettlementPrice    float64 //本次结算价
	UpperLimitPrice    float64 //涨停板价
	LowerLimitPrice    float64 //跌停板价
	PreDelta           float64 //昨虚实度
	CurrentDelta       float64 //今虚实度
	UpdateTime         string  //最后修改时间
	UpdateMillisec     string  //最后修改毫秒
	BidPrice1          float64 //申买价一
	BidVolume1         int64   //申买量一
	AskPrice1          float64 //申卖价一
	AskVolume1         int64   //申卖量一
	BidPrice2          float64 //申买价二
	BidVolume2         int64   //申买量二
	AskPrice2          float64 //申卖价二
	AskVolume2         int64   //申卖量二
	BidPrice3          float64 //申买价三
	BidVolume3         int64   //申买量三
	AskPrice3          float64 //申卖价三
	AskVolume3         int64   //申卖量三
	BidPrice4          float64 //申买价四
	BidVolume4         int64   //申买量四
	AskPrice4          float64 //申卖价四
	AskVolume4         int64   //申卖量四
	BidPrice5          float64 //申买价五
	BidVolume5         int64   //申买量五
	AskPrice5          float64 //申卖价五
	AskVolume5         int64   //申卖量五
	AveragePrice       float64 //当日均价
	XSpeedTime         string  //系统柜台时间
	TimeStamp          string  //成交时间戳
}

//K线数据
type CandleStickData struct {
	KNum           int64         //K线编号
	KDuration      time.Duration //K线周期
	InstrumentID   string        //合约代码
	KTime          time.Time     //时间K
	OpenPrice      float64       //开
	OpenTimestamp  time.Time     //第一条行情时间戳
	ClosePrice     float64       //收
	CloseTimestamp time.Time     //最后一条行情时间戳
	HighestPrice   float64       //最高
	LowestPrice    float64       //最低
	Volume         int64         //总手
	OpenInterest   int64         //持仓量
	Expansion      KExpansion    //K线扩展，复杂属性，调用GetComplexCandleStickData()函数才会被赋值
	Count          int64         //包含的行情条数
}

//K线扩展,复杂属性
type KExpansion struct {
	PriceVolumeTotal float64 //价格*量之和
	VolumeTotal      int64   //总量
	Vwap             float64 //成交量加权平均价
}

//K线周期和复杂度组合
type KParams struct {
	KComplexity int           //K线复杂度
	KDuration   time.Duration //K线周期
}

const (
	K_NONE    = 0 //不计算K线
	K_SIMPLE  = 1 //计算简单K线
	K_COMPLEX = 2 //计算复杂K线，包括：加权平均价
)

//持仓数据
type PositionData struct {
	AccountID              string  //账号
	InstrumentID           string  //合约号
	BuySellType            string  //买卖标志
	TotalAvaiAmount        int64   //总可用
	AveragePrice           float64 //平均成本
	DatePositionProfitLoss float64 //盯视持仓盈亏
	DateCloseProfitLoss    float64 //盯视平仓盈亏
	FloatProfitLoss        float64 //浮动盈亏
	IsFirst                int64   //是否是第一条
	IsLast                 int64   //是否是最后一条
}

//本地持仓数据
type LocalPositionData struct {
	AccountID              string  //账号
	InstrumentID           string  //合约号
	BuySellType            string  //买卖标志
	TotalAvaiAmount        int64   //总可用
	AveragePrice           float64 //平均成本
	DatePositionProfitLoss float64 //盯视持仓盈亏
	DateCloseProfitLoss    float64 //盯视平仓盈亏
	FloatProfitLoss        float64 //浮动盈亏
	Strategy               string  //策略实例名
}

//订单数据
type OrderData struct {
	InstrumentCode string  //合约代码
	InstrumentNum  string  //合约号
	BuySellType    string  //买卖标志 buy 买 sell 卖
	OpenCloseType  string  //开平标志 open 开 closetoday，closeyesterday 平仓
	OrderAmount    int64   //报单数量
	OrderPriceType string  //报价类型，limitprice 限价，marketprice 市价
	InsertPrice    float64 //报单价格
	LocalOrderID   int64   //本地下单id
	Strategy       string  //策略名称
}

//撤单数据
type CancelOrderData struct {
	InstrumentCode string //合约代码
	InstrumentNum  string //合约号
	LocalOrderID   int64  //本地下单id
}

//下单响应数据
type PlaceOrderResponseData struct {
	ResponseErrorID string //错误代码 0=成功,!0=不成功
	LocalOrderID    int64  //本地下单id
}

////委托回报数据
//type OrderPlacedReturnData struct {
//	StatusCode   string
//	StatusMsg    string
//	SysOrderID   string
//	ReturnTime   int64
//	LocalOrderID int64 //本地下单id
//}

//成交回报数据
type OrderMatchedReturnData struct {
	//InstrumentID  string  //合约号
	TradePrice   float64 //成交价格
	TradeAmount  int64   //成交量
	TradeTime    string  //交易所返回时间
	LocalOrderID int64   //本地下单id
	MatchedTime  int64   //程序时间戳
	//OpenCloseType string  //开平仓标志
}

//撤单响应数据
type CancelOrderResponseData struct {
	ResponseErrorID string //错误代码 0=成功,!0=不成功
	LocalOrderID    int64  //本地下单id
}

//撤单回报数据
type OrderCanceledReturnData struct {
	LocalOrderID  int64  //本地下单id
	CancelAmount  int64  //成交量
	CancelTime    string //交易所返回时间
	CancelSysTime int64  //程序时间戳
}

//资金数据
type CapitalData struct {
	Balance       float64 //动态权益
	CurrentMargin float64 //占用保证金
	Available     float64 //可用资金
}

//合约规则
type InstrumentRule struct {
	Code        string  //代码
	Name        string  //名称
	UnitAmount  int     //交易单位
	Months      []int   //合约月份
	MarginRatio float64 //风险率，保证金率
}

//MACD数据
type MACDData struct {
	EMAFast float64
	EMASlow float64
	Dif     float64
	Dea     float64 //DEA又称:DEM,MACD
}

//KD数据
type KDData struct {
	K float64
	D float64
}

//频道信息
const (
	CH_MARKET           = "STRATEGY:MarketData"           //行情频道
	CH_POSITION_PUBLISH = "STRATEGY:ReqQryPosition2"      //发送持仓查询的频道
	CH_POSITION         = "TRADE:LatestCompletedPosition" //返回持仓信息的频道

	CH_INSERT_ORDER             = "STRATEGY:ReqInsertOrder" //下单频道
	CH_INSERT_ORDER_RSP         = "TRADE:OnRspInsertOrder"  //下单响应频道
	CH_INSERT_ORDER_MATCHED_RTN = "TRADE:OnRtnMatchedInfo"  //成交回报频道

	CH_CANCEL_ORDER     = "STRATEGY:ReqCancelOrder" //撤单频道
	CH_CANCEL_ORDER_RSP = "TRADE:OnRspCancelOrder"  //撤单响应频道
	CH_CANCEL_ORDER_RTN = "TRADE:OnRtnCancelOrder"  //撤单回报频道

	CH_CAPITAL_PUBLISH = "STRATEGY:ReqQryCustomerCapital" //发送资金查询频道
	CH_CAPITAL         = "TRADE:TradingAccount"           //获取资金频道

	CH_TEST_MARKET        = "QUOTE:test"     //回测行情频道
	CH_TEST_GET_TIME_TICK = "STRATEGY:ReqMd" //通知获取下一个时间tick的频道
)

//Redis变量
const (
	VAR8_INST_RULES = "QUOTE:instrumentsInfo"
)

//订单状态
const (
	ORDER_STATUS_PLACED        = iota //刚刚下单
	ORDER_STATUS_PART_DONE            //订单部分成交
	ORDER_STATUS_DONE                 //订单完全成交
	ORDER_STATUS_PART_CANCELED        //订单部分撤销
	ORDER_STATUS_CANCELED             //订单完全撤销
)

//时间格式
const (
	YYYY_MM_DD_hh_mm_ss = "2006-01-02 15:04:05.999"
	YYYY_MM_DD          = "2006-01-02"
	hh_mm_ss            = "15:04:05"

	YYYY_MM_DD_hh_mm_ss_MST    = "2006-01-02 15:04:05.999 MST"
	YYYY_MM_DD_hh_mm_ss_offset = "2006-01-02 15:04:05.999 -0700"
	YYYYMMDDhhmmss_MST         = "20060102150405 MST"
)

//各大交易所交易时间
const (
	DAY_C1   = iota //无夜市
	NIGHT_C1        //有夜市23：00收盘
	NIGHT_C2        //有夜市23：30收盘
	NIGHT_C3        //有夜市01：00收盘
	NIGHT_C4        //有夜市02：30收盘
	ZJ_DAY          //中金所品种，无夜市
)

const (
	DAY_BID_START      = "08:55:00"
	DAY_BID_END        = "09:00:00"
	DAY_TRADE_START_P1 = "09:00:00"
	DAY_TRADE_END_P1   = "10:15:00"
	DAY_TRADE_START_P2 = "10:30:00"
	DAY_TRADE_END_P2   = "11:30:00"
	DAY_TRADE_START_P3 = "13:30:00"
	DAY_TRADE_END_P3   = "15:00:00"

	NIGHT_C1_BID_START      = "20:55:00"
	NIGHT_C1_BID_END        = "21:00:00"
	NIGHT_C1_TRADE_START_P1 = "21:00:00"
	NIGHT_C1_TRADE_END_P1   = "23:00:00"

	NIGHT_C2_BID_START      = "20:55:00"
	NIGHT_C2_BID_END        = "21:00:00"
	NIGHT_C2_TRADE_START_P1 = "21:00:00"
	NIGHT_C2_TRADE_END_P1   = "23:30:00"

	NIGHT_C3_BID_START      = "20:55:00"
	NIGHT_C3_BID_END        = "21:00:00"
	NIGHT_C3_TRADE_START_P1 = "21:00:00"
	NIGHT_C3_TRADE_END_P1   = "01:00:00" //次日

	NIGHT_C4_BID_START      = "20:55:00"
	NIGHT_C4_BID_END        = "21:00:00"
	NIGHT_C4_TRADE_START_P1 = "21:00:00"
	NIGHT_C4_TRADE_END_P1   = "02:30:00" //次日

	//中金所品种
	ZJ_DAY_BID_START      = "09:15:00"
	ZJ_DAY_BID_END        = "09:30:00"
	ZJ_DAY_TRADE_START_P1 = "09:30:00"
	ZJ_DAY_TRADE_END_P1   = "11:30:00"
	ZJ_DAY_TRADE_START_P2 = "13:00:00"
	ZJ_DAY_TRADE_END_P2   = "15:00:00"
)
