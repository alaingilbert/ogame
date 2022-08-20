package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
)

// LazyResearches ...
type LazyResearches func() Researches

func (s LazyResearches) GetEnergyTechnology() int64     { return s().EnergyTechnology }
func (s LazyResearches) GetLaserTechnology() int64      { return s().LaserTechnology }
func (s LazyResearches) GetIonTechnology() int64        { return s().IonTechnology }
func (s LazyResearches) GetHyperspaceTechnology() int64 { return s().HyperspaceTechnology }
func (s LazyResearches) GetPlasmaTechnology() int64     { return s().PlasmaTechnology }
func (s LazyResearches) GetCombustionDrive() int64      { return s().CombustionDrive }
func (s LazyResearches) GetImpulseDrive() int64         { return s().ImpulseDrive }
func (s LazyResearches) GetHyperspaceDrive() int64      { return s().HyperspaceDrive }
func (s LazyResearches) GetEspionageTechnology() int64  { return s().EspionageTechnology }
func (s LazyResearches) GetComputerTechnology() int64   { return s().ComputerTechnology }
func (s LazyResearches) GetAstrophysics() int64         { return s().Astrophysics }
func (s LazyResearches) GetIntergalacticResearchNetwork() int64 {
	return s().IntergalacticResearchNetwork
}
func (s LazyResearches) GetGravitonTechnology() int64  { return s().GravitonTechnology }
func (s LazyResearches) GetWeaponsTechnology() int64   { return s().WeaponsTechnology }
func (s LazyResearches) GetShieldingTechnology() int64 { return s().ShieldingTechnology }
func (s LazyResearches) GetArmourTechnology() int64    { return s().ArmourTechnology }

// Researches represent player's researches
type Researches struct {
	EnergyTechnology             int64 // 113
	LaserTechnology              int64 // 120
	IonTechnology                int64 // 121
	HyperspaceTechnology         int64 // 114
	PlasmaTechnology             int64 // 122
	CombustionDrive              int64 // 115
	ImpulseDrive                 int64 // 117
	HyperspaceDrive              int64 // 118
	EspionageTechnology          int64 // 106
	ComputerTechnology           int64 // 108
	Astrophysics                 int64 // 124
	IntergalacticResearchNetwork int64 // 123
	GravitonTechnology           int64 // 199
	WeaponsTechnology            int64 // 109
	ShieldingTechnology          int64 // 110
	ArmourTechnology             int64 // 111
}

func (s Researches) GetEnergyTechnology() int64             { return s.EnergyTechnology }
func (s Researches) GetLaserTechnology() int64              { return s.LaserTechnology }
func (s Researches) GetIonTechnology() int64                { return s.IonTechnology }
func (s Researches) GetHyperspaceTechnology() int64         { return s.HyperspaceTechnology }
func (s Researches) GetPlasmaTechnology() int64             { return s.PlasmaTechnology }
func (s Researches) GetCombustionDrive() int64              { return s.CombustionDrive }
func (s Researches) GetImpulseDrive() int64                 { return s.ImpulseDrive }
func (s Researches) GetHyperspaceDrive() int64              { return s.HyperspaceDrive }
func (s Researches) GetEspionageTechnology() int64          { return s.EspionageTechnology }
func (s Researches) GetComputerTechnology() int64           { return s.ComputerTechnology }
func (s Researches) GetAstrophysics() int64                 { return s.Astrophysics }
func (s Researches) GetIntergalacticResearchNetwork() int64 { return s.IntergalacticResearchNetwork }
func (s Researches) GetGravitonTechnology() int64           { return s.GravitonTechnology }
func (s Researches) GetWeaponsTechnology() int64            { return s.WeaponsTechnology }
func (s Researches) GetShieldingTechnology() int64          { return s.ShieldingTechnology }
func (s Researches) GetArmourTechnology() int64             { return s.ArmourTechnology }

// ToPtr returns a pointer to self
func (s Researches) ToPtr() *Researches {
	return &s
}

// Lazy returns a function that return self
func (s Researches) Lazy() LazyResearches {
	return func() Researches { return s }
}

// ByID gets the player research level by research id
func (s Researches) ByID(id ID) int64 {
	return researchByID(id, s)
}

func researchByID(id ID, researches IResearches) int64 {
	switch id {
	case EnergyTechnologyID:
		return researches.GetEnergyTechnology()
	case LaserTechnologyID:
		return researches.GetLaserTechnology()
	case IonTechnologyID:
		return researches.GetIonTechnology()
	case HyperspaceTechnologyID:
		return researches.GetHyperspaceTechnology()
	case PlasmaTechnologyID:
		return researches.GetPlasmaTechnology()
	case CombustionDriveID:
		return researches.GetCombustionDrive()
	case ImpulseDriveID:
		return researches.GetImpulseDrive()
	case HyperspaceDriveID:
		return researches.GetHyperspaceDrive()
	case EspionageTechnologyID:
		return researches.GetEspionageTechnology()
	case ComputerTechnologyID:
		return researches.GetComputerTechnology()
	case AstrophysicsID:
		return researches.GetAstrophysics()
	case IntergalacticResearchNetworkID:
		return researches.GetIntergalacticResearchNetwork()
	case GravitonTechnologyID:
		return researches.GetGravitonTechnology()
	case WeaponsTechnologyID:
		return researches.GetWeaponsTechnology()
	case ShieldingTechnologyID:
		return researches.GetShieldingTechnology()
	case ArmourTechnologyID:
		return researches.GetArmourTechnology()
	}
	return 0
}

func (s Researches) String() string {
	return "\n" +
		"             Energy Technology: " + utils.FI64(s.EnergyTechnology) + "\n" +
		"              Laser Technology: " + utils.FI64(s.LaserTechnology) + "\n" +
		"                Ion Technology: " + utils.FI64(s.IonTechnology) + "\n" +
		"         Hyperspace Technology: " + utils.FI64(s.HyperspaceTechnology) + "\n" +
		"             Plasma Technology: " + utils.FI64(s.PlasmaTechnology) + "\n" +
		"              Combustion Drive: " + utils.FI64(s.CombustionDrive) + "\n" +
		"                 Impulse Drive: " + utils.FI64(s.ImpulseDrive) + "\n" +
		"              Hyperspace Drive: " + utils.FI64(s.HyperspaceDrive) + "\n" +
		"          Espionage Technology: " + utils.FI64(s.EspionageTechnology) + "\n" +
		"           Computer Technology: " + utils.FI64(s.ComputerTechnology) + "\n" +
		"                  Astrophysics: " + utils.FI64(s.Astrophysics) + "\n" +
		"Intergalactic Research Network: " + utils.FI64(s.IntergalacticResearchNetwork) + "\n" +
		"           Graviton Technology: " + utils.FI64(s.GravitonTechnology) + "\n" +
		"            Weapons Technology: " + utils.FI64(s.WeaponsTechnology) + "\n" +
		"          Shielding Technology: " + utils.FI64(s.ShieldingTechnology) + "\n" +
		"             Armour Technology: " + utils.FI64(s.ArmourTechnology)
}
