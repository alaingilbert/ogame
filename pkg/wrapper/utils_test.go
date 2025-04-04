package wrapper

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertIntoCoordinate(t *testing.T) {
	expected := ogame.MustParseCoord("1:1:1")
	assert.Equal(t, expected, utils.First(ConvertIntoCoordinate(nil, "1:1:1")))
	assert.Equal(t, expected, utils.First(ConvertIntoCoordinate(nil, ogame.MustParseCoord("1:1:1"))))
	assert.Equal(t, expected, utils.First(ConvertIntoCoordinate(nil, ogame.EmpireCelestial{Coordinate: ogame.MustParseCoord("1:1:1")})))
	assert.Equal(t, expected, utils.First(ConvertIntoCoordinate(nil, ogame.Planet{Coordinate: ogame.MustParseCoord("1:1:1")})))
	assert.Equal(t, expected, utils.First(ConvertIntoCoordinate(nil, ogame.Moon{Coordinate: ogame.MustParseCoord("1:1:1")})))
}
