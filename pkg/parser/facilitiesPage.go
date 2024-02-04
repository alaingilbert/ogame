package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *FacilitiesPage) ExtractFacilities() (ogame.Facilities, error) {
	return p.e.ExtractFacilitiesFromDoc(p.GetDoc())
}
