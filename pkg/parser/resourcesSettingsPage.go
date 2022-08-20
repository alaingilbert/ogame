package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p ResourcesSettingsPage) ExtractResourceSettings() (ogame.ResourceSettings, error) {
	return p.e.ExtractResourceSettingsFromDoc(p.GetDoc())
}
