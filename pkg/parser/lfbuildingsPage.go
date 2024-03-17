package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *LfBuildingsPage) ExtractLfBuildings() (ogame.LfBuildings, error) {
	return p.e.ExtractLfBuildingsFromDoc(p.GetDoc())
}
