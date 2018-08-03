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

// ByOGameID ...
func (s Researches) ByOGameID(ogameID ID) int {
	if ogameID == EnergyTechnology.ID {
		return s.EnergyTechnology
	} else if ogameID == LaserTechnology.ID {
		return s.LaserTechnology
	} else if ogameID == IonTechnology.ID {
		return s.IonTechnology
	} else if ogameID == HyperspaceTechnology.ID {
		return s.HyperspaceTechnology
	} else if ogameID == PlasmaTechnology.ID {
		return s.PlasmaTechnology
	} else if ogameID == CombustionDrive.ID {
		return s.CombustionDrive
	} else if ogameID == ImpulseDrive.ID {
		return s.ImpulseDrive
	} else if ogameID == HyperspaceDrive.ID {
		return s.HyperspaceDrive
	} else if ogameID == EspionageTechnology.ID {
		return s.EspionageTechnology
	} else if ogameID == ComputerTechnology.ID {
		return s.ComputerTechnology
	} else if ogameID == Astrophysics.ID {
		return s.Astrophysics
	} else if ogameID == IntergalacticResearchNetwork.ID {
		return s.IntergalacticResearchNetwork
	} else if ogameID == GravitonTechnology.ID {
		return s.GravitonTechnology
	} else if ogameID == WeaponsTechnology.ID {
		return s.WeaponsTechnology
	} else if ogameID == ShieldingTechnology.ID {
		return s.ShieldingTechnology
	} else if ogameID == ArmourTechnology.ID {
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
