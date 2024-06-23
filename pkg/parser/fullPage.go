package parser

import (
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
)

func (p *FullPage) ExtractOGameSession() string {
	return p.e.ExtractOGameSessionFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractIsInVacation() bool {
	return p.e.ExtractIsInVacationFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractAjaxChatToken() (string, error) {
	return p.e.ExtractAjaxChatToken(p.content)
}

func (p *FullPage) ExtractToken() (string, error) {
	return p.e.ExtractToken(p.content)
}

func (p *FullPage) ExtractCharacterClass() (ogame.CharacterClass, error) {
	return p.e.ExtractCharacterClassFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractCommander() bool {
	return p.e.ExtractCommanderFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractAdmiral() bool {
	return p.e.ExtractAdmiralFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractLifeformTypeFromDoc() ogame.LifeformType {
	return p.e.ExtractLifeformTypeFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractEngineer() bool {
	return p.e.ExtractEngineerFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractGeologist() bool {
	return p.e.ExtractGeologistFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractTechnocrat() bool {
	return p.e.ExtractTechnocratFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractColonies() (int64, int64) {
	return p.e.ExtractColoniesFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractLifeformEnabled() bool {
	return p.e.ExtractLifeformEnabled(p.GetContent())
}

func (p *FullPage) ExtractServerTime() (time.Time, error) {
	return p.e.ExtractServerTimeFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractPlanets() []ogame.Planet {
	return p.e.ExtractPlanetsFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractPlanet(v any) (ogame.Planet, error) {
	return p.e.ExtractPlanetFromDoc(p.GetDoc(), v)
}

func (p *FullPage) ExtractPlanetCoordinate() (ogame.Coordinate, error) {
	return p.e.ExtractPlanetCoordinate(p.content)
}

func (p *FullPage) ExtractPlanetID() (ogame.CelestialID, error) {
	return p.e.ExtractPlanetID(p.content)
}

func (p *FullPage) ExtractMoons() []ogame.Moon {
	return p.e.ExtractMoonsFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractMoon(v any) (ogame.Moon, error) {
	return p.e.ExtractMoonFromDoc(p.GetDoc(), v)
}

func (p *FullPage) ExtractCelestials() ([]ogame.Celestial, error) {
	return p.e.ExtractCelestialsFromDoc(p.GetDoc())
}

func (p *FullPage) ExtractCelestial(v any) (ogame.Celestial, error) {
	return p.e.ExtractCelestialFromDoc(p.GetDoc(), v)
}

func (p *FullPage) ExtractResources() ogame.Resources {
	return p.e.ExtractResourcesFromDoc(p.GetDoc())
}
