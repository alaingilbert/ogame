package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *ResearchPage) ExtractResearch() ogame.Researches {
	return p.e.ExtractResearchFromDoc(p.GetDoc())
}
