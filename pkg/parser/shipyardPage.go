package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *ShipyardPage) ExtractProduction() ([]ogame.Quantifiable, int64, error) {
	return p.e.ExtractProduction(p.content)
}

func (p *ShipyardPage) ExtractShips() (ogame.ShipsInfos, error) {
	return p.e.ExtractShipsFromDoc(p.GetDoc())
}
