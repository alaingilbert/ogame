package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p SuppliesPage) ExtractLfBuildings() (ogame.LfBuildings, error) {
	return p.e.ExtractLfBuildingsFromDoc(p.GetDoc())
}
