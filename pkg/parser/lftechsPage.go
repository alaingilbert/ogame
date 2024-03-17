package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *LfResearchPage) ExtractLfResearch() (ogame.LfResearches, error) {
	return p.e.ExtractLfResearchFromDoc(p.GetDoc())
}
