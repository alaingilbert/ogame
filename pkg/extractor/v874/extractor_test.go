package v874

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExtractOfferOfTheDayPrice(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v8.7.4/en/traderImportExport.html")
	price, token, _, _, _ := NewExtractor().ExtractOfferOfTheDay(pageHTMLBytes)
	assert.Equal(t, int64(178224), price)
	assert.Equal(t, "2a38193e2fa6047e1d92d2f2c71c00fd", token)
}

func TestExtractAuction(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v8.7.4/en/traderAuctioneer.html")
	res, _ := NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, "43576386810cdf91a833a6239f323f66", res.Token)
}
