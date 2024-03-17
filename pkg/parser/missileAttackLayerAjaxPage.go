package parser

func (p *MissileAttackLayerAjaxPage) ExtractIPM() (int64, int64, string) {
	return p.e.ExtractIPMFromDoc(p.GetDoc())
}
