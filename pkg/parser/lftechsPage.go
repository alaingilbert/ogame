package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p LfTechsPage) ExtractLfTechs() (ogame.LfTechs, error) {
	return p.e.ExtractLfTechsFromDoc(p.GetDoc())
}
