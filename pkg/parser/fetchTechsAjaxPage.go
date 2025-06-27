package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *FetchTechsAjaxPage) ExtractTechs() (ogame.Techs, error) {
	return p.e.ExtractTechs(p.content)
}
