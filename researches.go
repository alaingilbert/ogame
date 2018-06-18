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
	if ogameID == EnergyTechnology {
		return s.EnergyTechnology
	} else if ogameID == LaserTechnology {
		return s.LaserTechnology
	} else if ogameID == IonTechnology {
		return s.IonTechnology
	} else if ogameID == HyperspaceTechnology {
		return s.HyperspaceTechnology
	} else if ogameID == PlasmaTechnology {
		return s.PlasmaTechnology
	} else if ogameID == CombustionDrive {
		return s.CombustionDrive
	} else if ogameID == ImpulseDrive {
		return s.ImpulseDrive
	} else if ogameID == HyperspaceDrive {
		return s.HyperspaceDrive
	} else if ogameID == EspionageTechnology {
		return s.EspionageTechnology
	} else if ogameID == ComputerTechnology {
		return s.ComputerTechnology
	} else if ogameID == Astrophysics {
		return s.Astrophysics
	} else if ogameID == IntergalacticResearchNetwork {
		return s.IntergalacticResearchNetwork
	} else if ogameID == GravitonTechnology {
		return s.GravitonTechnology
	} else if ogameID == WeaponsTechnology {
		return s.WeaponsTechnology
	} else if ogameID == ShieldingTechnology {
		return s.ShieldingTechnology
	} else if ogameID == ArmourTechnology {
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
