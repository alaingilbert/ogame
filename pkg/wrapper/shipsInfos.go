package wrapper

import "github.com/alaingilbert/ogame/pkg/ogame"

// ShipsInfos ...
type ShipsInfos struct {
	wrapper Wrapper
	ogame.ShipsInfos
}

// NewShipsInfos ...
func NewShipsInfos(w Wrapper, shipsInfos ogame.ShipsInfos) *ShipsInfos {
	return &ShipsInfos{
		wrapper:    w,
		ShipsInfos: shipsInfos,
	}
}

// Cargo ...
func (s *ShipsInfos) Cargo() (out int64) {
	tx := s.wrapper.Begin()
	defer tx.Done()
	techs := tx.GetCachedResearch()
	bonuses, _ := tx.GetCachedLfBonuses()
	characterClass := s.wrapper.CharacterClass()
	multiplier := float64(s.wrapper.GetServerData().CargoHyperspaceTechMultiplier) / 100
	probeRaids := s.wrapper.GetServer().Settings.EspionageProbeRaids == 1
	return s.ShipsInfos.Cargo(techs, bonuses.LfShipBonuses, characterClass, multiplier, probeRaids)
}

// Speed ...
func (s *ShipsInfos) Speed() (out int64) {
	tx := s.wrapper.Begin()
	defer tx.Done()
	techs := tx.GetCachedResearch()
	bonuses, _ := tx.GetCachedLfBonuses()
	characterClass := s.wrapper.CharacterClass()
	return s.ShipsInfos.Speed(techs, bonuses.LfShipBonuses, characterClass)
}
