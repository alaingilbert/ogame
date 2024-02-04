package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *PhalanxAjaxPage) ExtractPhalanx() ([]ogame.Fleet, error) {
	return p.e.ExtractPhalanx(p.content)
}
