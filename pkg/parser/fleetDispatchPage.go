package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *FleetDispatchPage) ExtractShips() (ogame.ShipsInfos, error) {
	return p.e.ExtractShipsFromDoc(p.GetDoc())
}

func (p *FleetDispatchPage) ExtractSlots() (ogame.Slots, error) {
	return p.e.ExtractSlotsFromDoc(p.GetDoc())
}

func (p *FleetDispatchPage) ExtractAcsValues() []ogame.ACSValues {
	return p.e.ExtractFleetDispatchACSFromDoc(p.GetDoc())
}
