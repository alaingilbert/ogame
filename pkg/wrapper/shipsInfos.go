package wrapper

import "github.com/alaingilbert/ogame/pkg/ogame"

type IShipsInfos interface {
	ByID(id ogame.ID) int64
	Cargo(techs ogame.IResearches, lfBonuses ogame.LfBonuses, characterClass ogame.CharacterClass, multiplier float64, probeRaids bool) int64
	Set(ogame.ID, int64)
	ToQuantifiables() []ogame.Quantifiable
}

// ShipsInfos ...
type ShipsInfos struct {
	wrapper Wrapper
	ogame.ShipsInfos
}

func (s *ShipsInfos) Cargo(techs ogame.IResearches, lfBonuses ogame.LfBonuses, characterClass ogame.CharacterClass, multiplier float64, probeRaids bool) int64 {
	return s.ShipsInfos.Cargo(techs, lfBonuses, characterClass, multiplier, probeRaids)
}
func (s *ShipsInfos) ByID(id ogame.ID) int64                { return s.ShipsInfos.ByID(id) }
func (s *ShipsInfos) Set(id ogame.ID, nbr int64)            { s.ShipsInfos.Set(id, nbr) }
func (s *ShipsInfos) ToQuantifiables() []ogame.Quantifiable { return s.ShipsInfos.ToQuantifiables() }

// NewShipsInfos ...
func NewShipsInfos(w Wrapper, shipsInfos ogame.ShipsInfos) *ShipsInfos {
	return &ShipsInfos{
		wrapper:    w,
		ShipsInfos: shipsInfos,
	}
}

// GetCargo ...
func (s *ShipsInfos) GetCargo() (out int64) {
	tx := s.wrapper.Begin()
	defer tx.Done()
	techs := tx.GetCachedResearch()
	bonuses, _ := tx.GetCachedLfBonuses()
	characterClass := s.wrapper.CharacterClass()
	multiplier := float64(s.wrapper.GetServerData().CargoHyperspaceTechMultiplier) / 100
	probeRaids := s.wrapper.GetServer().Settings.EspionageProbeRaids == 1
	return s.ShipsInfos.Cargo(techs, bonuses, characterClass, multiplier, probeRaids)
}

// Speed ...
func (s *ShipsInfos) Speed() (out int64) {
	tx := s.wrapper.Begin()
	defer tx.Done()
	techs := tx.GetCachedResearch()
	bonuses, _ := tx.GetCachedLfBonuses()
	allianceClass, _ := tx.GetCachedAllianceClass()
	characterClass := s.wrapper.CharacterClass()
	return s.ShipsInfos.Speed(techs, bonuses, characterClass, allianceClass)
}
