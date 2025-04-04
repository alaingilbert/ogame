package parser

func (p *MissileAttackLayerAjaxPage) ExtractIPM() (int64, int64, string, error) {
	return p.e.ExtractIPMFromDoc(p.GetDoc())
}
