package ogame

import "strconv"

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
	if id == EnergyTechnology.ID {
		return s.EnergyTechnology
	} else if id == LaserTechnology.ID {
		return s.LaserTechnology
	} else if id == IonTechnology.ID {
		return s.IonTechnology
	} else if id == HyperspaceTechnology.ID {
		return s.HyperspaceTechnology
	} else if id == PlasmaTechnology.ID {
		return s.PlasmaTechnology
	} else if id == CombustionDrive.ID {
		return s.CombustionDrive
	} else if id == ImpulseDrive.ID {
		return s.ImpulseDrive
	} else if id == HyperspaceDrive.ID {
		return s.HyperspaceDrive
	} else if id == EspionageTechnology.ID {
		return s.EspionageTechnology
	} else if id == ComputerTechnology.ID {
		return s.ComputerTechnology
	} else if id == Astrophysics.ID {
		return s.Astrophysics
	} else if id == IntergalacticResearchNetwork.ID {
		return s.IntergalacticResearchNetwork
	} else if id == GravitonTechnology.ID {
		return s.GravitonTechnology
	} else if id == WeaponsTechnology.ID {
		return s.WeaponsTechnology
	} else if id == ShieldingTechnology.ID {
		return s.ShieldingTechnology
	} else if id == ArmourTechnology.ID {
		return s.ArmourTechnology
	}
	return 0
}

func (s Researches) String() string {
	return "\n" +
		"             Energy Technology: " + strconv.FormatInt(s.EnergyTechnology, 10) + "\n" +
		"              Laser Technology: " + strconv.FormatInt(s.LaserTechnology, 10) + "\n" +
		"                Ion Technology: " + strconv.FormatInt(s.IonTechnology, 10) + "\n" +
		"         Hyperspace Technology: " + strconv.FormatInt(s.HyperspaceTechnology, 10) + "\n" +
		"             Plasma Technology: " + strconv.FormatInt(s.PlasmaTechnology, 10) + "\n" +
		"              Combustion Drive: " + strconv.FormatInt(s.CombustionDrive, 10) + "\n" +
		"                 Impulse Drive: " + strconv.FormatInt(s.ImpulseDrive, 10) + "\n" +
		"              Hyperspace Drive: " + strconv.FormatInt(s.HyperspaceDrive, 10) + "\n" +
		"          Espionage Technology: " + strconv.FormatInt(s.EspionageTechnology, 10) + "\n" +
		"           Computer Technology: " + strconv.FormatInt(s.ComputerTechnology, 10) + "\n" +
		"                  Astrophysics: " + strconv.FormatInt(s.Astrophysics, 10) + "\n" +
		"Intergalactic Research Network: " + strconv.FormatInt(s.IntergalacticResearchNetwork, 10) + "\n" +
		"           Graviton Technology: " + strconv.FormatInt(s.GravitonTechnology, 10) + "\n" +
		"            Weapons Technology: " + strconv.FormatInt(s.WeaponsTechnology, 10) + "\n" +
		"          Shielding Technology: " + strconv.FormatInt(s.ShieldingTechnology, 10) + "\n" +
		"             Armour Technology: " + strconv.FormatInt(s.ArmourTechnology, 10)
}
