package v11_13_0

import (
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_9_0"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// Extractor ...
type Extractor struct {
	v11_9_0.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractConstructions ...
func (e *Extractor) ExtractConstructions(pageHTML []byte) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64, lfBuildingID ogame.ID, lfBuildingCountdown int64, lfResearchID ogame.ID, lfResearchCountdown int64) {
	return ExtractConstructions(pageHTML, clockwork.NewRealClock())
}
