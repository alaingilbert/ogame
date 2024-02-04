package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *DefensesPage) ExtractDefense() (ogame.DefensesInfos, error) {
	return p.e.ExtractDefenseFromDoc(p.GetDoc())
}
