package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
)

// LazyResearches ...
type LazyResearches func() Researches

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
	if id == EnergyTechnologyID {
		return s.EnergyTechnology
	} else if id == LaserTechnologyID {
		return s.LaserTechnology
	} else if id == IonTechnologyID {
		return s.IonTechnology
	} else if id == HyperspaceTechnologyID {
		return s.HyperspaceTechnology
	} else if id == PlasmaTechnologyID {
		return s.PlasmaTechnology
	} else if id == CombustionDriveID {
		return s.CombustionDrive
	} else if id == ImpulseDriveID {
		return s.ImpulseDrive
	} else if id == HyperspaceDriveID {
		return s.HyperspaceDrive
	} else if id == EspionageTechnologyID {
		return s.EspionageTechnology
	} else if id == ComputerTechnologyID {
		return s.ComputerTechnology
	} else if id == AstrophysicsID {
		return s.Astrophysics
	} else if id == IntergalacticResearchNetworkID {
		return s.IntergalacticResearchNetwork
	} else if id == GravitonTechnologyID {
		return s.GravitonTechnology
	} else if id == WeaponsTechnologyID {
		return s.WeaponsTechnology
	} else if id == ShieldingTechnologyID {
		return s.ShieldingTechnology
	} else if id == ArmourTechnologyID {
		return s.ArmourTechnology
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
