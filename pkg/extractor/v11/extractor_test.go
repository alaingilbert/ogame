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
