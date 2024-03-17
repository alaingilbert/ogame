package parser

import "github.com/alaingilbert/ogame/pkg/ogame"

func (p *PreferencesPage) ExtractPreferences() ogame.Preferences {
	return p.e.ExtractPreferencesFromDoc(p.GetDoc())
}
