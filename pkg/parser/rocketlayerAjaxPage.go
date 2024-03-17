package parser

func (p *RocketlayerAjaxPage) ExtractDestroyRockets() (int64, int64, string, error) {
	return p.e.ExtractDestroyRockets(p.content)
}
