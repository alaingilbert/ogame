package parser

import (
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	v71 "github.com/alaingilbert/ogame/pkg/extractor/v71"
	v9 "github.com/alaingilbert/ogame/pkg/extractor/v9"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExtractCelestial(t *testing.T) {
	pageHTML, _ := os.ReadFile("../../samples/v9.0.2/en/lifeform/overview_all_queues.html")
	p := FullPage{Page: Page{e: v9.NewExtractor(), content: pageHTML}}
	celestial, err := p.ExtractCelestial(33640820)
	assert.NoError(t, err)
	assert.Equal(t, ogame.CelestialID(33640820), celestial.GetID())
	celestial1, err := p.ExtractCelestial(celestial)
	assert.NoError(t, err)
	assert.Equal(t, ogame.CelestialID(33640820), celestial1.GetID())

	pageHTML, _ = os.ReadFile("../../samples/v7.1/en/moon_facilities.html")
	p = FullPage{Page: Page{e: v71.NewExtractor(), content: pageHTML}}
	celestial2, err := p.ExtractCelestial(33780773)
	assert.NoError(t, err)
	assert.Equal(t, ogame.CelestialID(33780773), celestial2.GetID())
	celestial3, err := p.ExtractCelestial(celestial2)
	assert.NoError(t, err)
	assert.Equal(t, ogame.CelestialID(33780773), celestial3.GetID())

	_, err = p.ExtractCelestial(Page{})
	assert.EqualError(t, err, v6.ErrUnsupportedType.Error())
}
