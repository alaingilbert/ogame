package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *MovementPage) ExtractFleets() []ogame.Fleet {
	return p.e.ExtractFleetsFromDoc(p.GetDoc())
}

func (p *MovementPage) ExtractSlots() (ogame.Slots, error) {
	return p.e.ExtractSlotsFromDoc(p.GetDoc())
}

func (p *MovementPage) ExtractCancelFleetToken(fleetID ogame.FleetID) (string, error) {
	return p.e.ExtractCancelFleetToken(p.content, fleetID)
}
