package ogame

import "strconv"

// Researches ...
type Researches struct {
	EnergyTechnology             int
	LaserTechnology              int
	IonTechnology                int
	HyperspaceTechnology         int
	PlasmaTechnology             int
	CombustionDrive              int
	ImpulseDrive                 int
	HyperspaceDrive              int
	EspionageTechnology          int
	ComputerTechnology           int
	Astrophysics                 int
	IntergalacticResearchNetwork int
	GravitonTechnology           int
	WeaponsTechnology            int
	ShieldingTechnology          int
	ArmourTechnology             int
}

// ByID ...
func (s Researches) ByID(id ID) int {
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
		"             Energy Technology: " + strconv.Itoa(s.EnergyTechnology) + "\n" +
		"              Laser Technology: " + strconv.Itoa(s.LaserTechnology) + "\n" +
		"                Ion Technology: " + strconv.Itoa(s.IonTechnology) + "\n" +
		"         Hyperspace Technology: " + strconv.Itoa(s.HyperspaceTechnology) + "\n" +
		"             Plasma Technology: " + strconv.Itoa(s.PlasmaTechnology) + "\n" +
		"              Combustion Drive: " + strconv.Itoa(s.CombustionDrive) + "\n" +
		"                 Impulse Drive: " + strconv.Itoa(s.ImpulseDrive) + "\n" +
		"              Hyperspace Drive: " + strconv.Itoa(s.HyperspaceDrive) + "\n" +
		"          Espionage Technology: " + strconv.Itoa(s.EspionageTechnology) + "\n" +
		"           Computer Technology: " + strconv.Itoa(s.ComputerTechnology) + "\n" +
		"                  Astrophysics: " + strconv.Itoa(s.Astrophysics) + "\n" +
		"Intergalactic Research Network: " + strconv.Itoa(s.IntergalacticResearchNetwork) + "\n" +
		"           Graviton Technology: " + strconv.Itoa(s.GravitonTechnology) + "\n" +
		"            Weapons Technology: " + strconv.Itoa(s.WeaponsTechnology) + "\n" +
		"          Shielding Technology: " + strconv.Itoa(s.ShieldingTechnology) + "\n" +
		"             Armour Technology: " + strconv.Itoa(s.ArmourTechnology)
}
