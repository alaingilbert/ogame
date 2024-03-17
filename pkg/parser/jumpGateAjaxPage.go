package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *JumpGateAjaxPage) ExtractJumpGate() (ogame.ShipsInfos, string, []ogame.MoonID, int64) {
	return p.e.ExtractJumpGate(p.content)
}
