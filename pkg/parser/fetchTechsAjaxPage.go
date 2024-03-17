package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *FetchTechsAjaxPage) ExtractTechs() (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, ogame.LfResearches, error) {
	return p.e.ExtractTechs(p.content)
}
