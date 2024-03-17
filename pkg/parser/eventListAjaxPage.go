package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *EventListAjaxPage) ExtractAttacks(ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	return p.e.ExtractAttacksFromDoc(p.GetDoc(), ownCoords)
}
