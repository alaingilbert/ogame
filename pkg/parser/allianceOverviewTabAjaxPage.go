package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

type AllianceOverviewTabRes struct {
	Target  string `json:"target"`
	Content struct {
		AllianceAllianceOverview string `json:"alliance/alliance_overview"`
	} `json:"content"`
	Files struct {
		Js  []string `json:"js"`
		CSS []string `json:"css"`
	} `json:"files"`
	Page struct {
		StateObj string `json:"stateObj"`
		Title    string `json:"title"`
		URL      string `json:"url"`
	} `json:"page"`
	ServerTime   int    `json:"serverTime"`
	NewAjaxToken string `json:"newAjaxToken"`
}

func (p *AllianceOverviewTabAjaxPage) ExtractAllianceClass() (ogame.AllianceClass, error) {
	return p.e.ExtractAllianceClass(p.content)
}
