// 将读原始数据文件，抽取数据的功能放在这个文件
//
package lev0

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// csv格式的数据文件，一行代表一日股票数据.
// 文件第一行为表头, 共18列，依次是：
//
// date,code,open,high,low,close,preclose,volume,amount,
// adjustflag,turn,tradestatus,
// pctChg,peTTM,pbMRQ,psTTM,pcfNcfTTM,isST
//
// 表头翻译对照如下：
// 日期，代码，开盘价，最高价，最低价，收盘价，前收盘价，成交量(股)，成交额(元)，
// 复权状态(1: 后复权，2: 前复权，3: 不复权)，换手率，交易状态(1: 正常，0: 停牌)，
// 涨跌幅(百分比)，动态市盈率，市净率，动态市销率，动态市现率，是否ST股(1: 是，0:否)
//
// 当股票停牌时，开盘价、最高价、最低价、收盘价都为前一日的收盘价，
// 但是将成交量、成交额记为0，换手率turn为空，读取数据时需要注意。

// 一般交易日前收盘价preClose等于昨日收盘价，
// 但是除权登记日，前收盘价根据股权登记日收盘价与分红、配股情况重新计算而得，
// 方法如下：
// 1、计算出息价：
// 		出息价 = 股息登记日的收盘价 - 每股所分红现金额
// 2、计算除权价
//		送股后的除权价 = 股权登记日的收盘价/(1 + 每股送红股数)
// 3、计算除权除息价
//		除权除息价 = (股权登记日的收盘价 - 每股分红现金 + 配股价*每股配股数) / (1 + 每股送股数 + 每股配股数)
// 前收盘价由交易所计算并公布。首发日的前收盘价等于首发价格
//
// 根据上面信息设计日股票数据结构：DayTickInfo
// 日期与股票代码设计为time.Time, string类型外，其余设计成数组以方便从文件读取数据
type Data struct {
	StockCode   string
	Dates       []time.Time
	Opens       []float64
	Highs       []float64
	Lows        []float64
	Closes      []float64
	PreCloses   []float64
	Volumns     []float64
	Amounts     []float64
	AdjustFlags []float64
	Turns       []float64
	TradeStatus []float64
	PctChgs     []float64
	PeTTMs      []float64
	PbMRQs      []float64
	PsTTMs      []float64
	PcfNctTTMs  []float64
	IsSTs       []float64
}

// ReadDayTickTable 从给定的数据文件 path 读取数据到 DataTable
// 并返回其指针。
func ReadData(path string) *Data {
	var (
		date time.Time
		code string
		ar   [16]float64
		errs [16]error
		err  error

		tbl Data
	)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	all, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	for i, row := range all[1:] {
		date, _ = time.Parse("2006-01-02", row[0])
		code = row[1]
		for j, s := range row[2:] {
			ar[j], errs[j] = strconv.ParseFloat(s, 64)
		}
		for j, err := range errs {
			if err != nil {
				if j == 8 {
					continue
				}
				fmt.Printf("err at (%d, %d)\n", i, j)
				fmt.Printf("parsed value = %f\n", ar[j])
			}
		}
		tbl.Dates = append(tbl.Dates, date)
		tbl.Opens = append(tbl.Opens, ar[0])
		tbl.Highs = append(tbl.Highs, ar[1])
		tbl.Lows = append(tbl.Lows, ar[2])
		tbl.Closes = append(tbl.Closes, ar[3])
		tbl.PreCloses = append(tbl.PreCloses, ar[4])
		tbl.Volumns = append(tbl.Volumns, ar[5])
		tbl.Amounts = append(tbl.Amounts, ar[6])
		tbl.AdjustFlags = append(tbl.AdjustFlags, ar[7])
		tbl.Turns = append(tbl.Turns, ar[8])
		tbl.TradeStatus = append(tbl.TradeStatus, ar[9])
		tbl.PctChgs = append(tbl.PctChgs, ar[10])
		tbl.PeTTMs = append(tbl.PeTTMs, ar[11])
		tbl.PbMRQs = append(tbl.PbMRQs, ar[12])
		tbl.PsTTMs = append(tbl.PsTTMs, ar[13])
		tbl.PcfNctTTMs = append(tbl.PcfNctTTMs, ar[14])
		tbl.IsSTs = append(tbl.IsSTs, ar[15])
	}
	tbl.StockCode = strings.Split(code, ".")[1]

	return &tbl
}
