package repository

import (
	"garantex/internal/models"
	util "garantex/internal/pkg/utils"
	"strconv"
)

const (
	table = "prices"
)

var (
	bufferPool = util.NewBufferPool[byte](1024)
)

func upsertPrice(price models.Price) string {
	buffer := bufferPool.Get()
	defer bufferPool.Put(buffer)
	buf := buffer.Data[:0]
	buf = append(buf, "INSERT INTO "...)
	buf = append(buf, table...)
	buf = append(buf, " (ts, ask_price, bid_price) VALUES "...)
	buf = append(buf, "("...)
	buf = append(buf, strconv.Itoa(int(price.Timestamp))...)
	buf = append(buf, ',')
	buf = strconv.AppendFloat(buf, price.AskPrice, 'f', -1, 32)
	buf = append(buf, ',')
	buf = strconv.AppendFloat(buf, price.BidPrice, 'f', -1, 32)
	buf = append(buf, ")"...)
	buf = append(buf, " ON CONFLICT (ts) DO UPDATE SET ts = EXCLUDED.ts, ask_price = EXCLUDED.ask_price, bid_price = EXCLUDED.bid_price;"...)
	return string(buf)
}
