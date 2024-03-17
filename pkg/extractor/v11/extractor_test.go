package v11

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExtractGalaxyInfos(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.6.2/galaxy1.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Gemini", 123, 456)
	assert.Equal(t, "Admiral Dorado", infos.Position(4).Player.Name)
	assert.Equal(t, "Commodore Gemini", infos.Position(10).Player.Name)
}

func TestCancel(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.6.2/overview_cancels.html")
	token, techID, listID, err := NewExtractor().ExtractCancelBuildingInfos(pageHTMLBytes)
	assert.Nil(t, err)
	assert.Equal(t, "5175965b16d3e743a710b8a07e5b35f1", token)
	assert.Equal(t, int64(1), techID)
	assert.Equal(t, int64(5168837), listID)
}
