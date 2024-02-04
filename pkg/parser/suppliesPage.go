package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *SuppliesPage) ExtractResourcesBuildings() (ogame.ResourcesBuildings, error) {
	return p.e.ExtractResourcesBuildingsFromDoc(p.GetDoc())
}
