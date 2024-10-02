package repository

import (
	"garantex/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expected = "INSERT INTO prices (ts, ask_price, bid_price) VALUES (1727834861,194.677,0.78415) ON CONFLICT (ts) DO UPDATE SET ts = EXCLUDED.ts, ask_price = EXCLUDED.ask_price, bid_price = EXCLUDED.bid_price;"

func TestUpsertPrice(t *testing.T) {
	price := models.Price{
		Timestamp: 1727834861,
		AskPrice:  194.677,
		BidPrice:  0.78415,
	}
	assert.Equal(t, expected, upsertPrice(price))
}
