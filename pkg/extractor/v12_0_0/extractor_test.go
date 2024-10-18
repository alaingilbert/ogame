package v12_0_0

import (
	"github.com/alaingilbert/clockwork"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestExtractServerTime(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v12.0.0/en/overview.html")
	clock := clockwork.NewFakeClockAt(time.Date(2024, 10, 18, 6, 20, 12, 0, time.UTC))
	res, err := NewExtractor().extractServerTime(pageHTMLBytes, clock)
	assert.Nil(t, err)
	assert.Equal(t, "2024-10-18 07:20:12 +0100 OGT", res.String())
}
