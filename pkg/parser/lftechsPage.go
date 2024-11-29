package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *LfResearchPage) ExtractLfResearch() (ogame.LfResearches, error) {
	return p.e.ExtractLfResearchFromDoc(p.GetDoc())
}

func (p *LfResearchPage) ExtractLfSlots() [18]ogame.LfSlot {
	return p.e.ExtractLfSlotsFromDoc(p.GetDoc())
}

func (p *LfResearchPage) ExtractArtefacts() (int64, int64) {
	return p.e.ExtractArtefactsFromDoc(p.GetDoc())
}
