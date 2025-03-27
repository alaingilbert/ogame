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

func TestExtractHighscore(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v12.0.27/en/highscore.html")
	highscore, _ := NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, int64(107694), highscore.Players[0].ID)
}

func TestExtractHighscoreIgnored(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v12.0.27/en/highscore_ignored.html")
	highscore, _ := NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, "Kosmokratoras", highscore.Players[66].Name)
	assert.Equal(t, int64(114011), highscore.Players[66].ID)
	assert.Equal(t, "Marshal Tempo", highscore.Players[67].Name)
	assert.Equal(t, int64(111580), highscore.Players[67].ID)
}

func TestExtractGalaxyInfos_alliance(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v12.0.27/en/galaxy_ajax.json")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(500077), infos.Position(9).Alliance.ID)
	assert.Equal(t, "NbE", infos.Position(9).Alliance.Tag)
	assert.Equal(t, int64(1), infos.Position(9).Alliance.Rank)
	assert.Equal(t, int64(110), infos.Position(9).Alliance.Member)
}
