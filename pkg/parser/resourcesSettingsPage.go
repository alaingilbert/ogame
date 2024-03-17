package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *ResourcesSettingsPage) ExtractResourceSettings() (ogame.ResourceSettings, string, error) {
	return p.e.ExtractResourceSettingsFromDoc(p.GetDoc())
}
