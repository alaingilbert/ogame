package ogame

import "strconv"

// Facilities ...
type Facilities struct {
	RoboticsFactory int
	Shipyard        int
	ResearchLab     int
	AllianceDepot   int
	MissileSilo     int
	NaniteFactory   int
	Terraformer     int
	SpaceDock       int
}

// ByOGameID ...
func (f Facilities) ByOGameID(ogameID ID) int {
	if ogameID == RoboticsFactory.ID {
		return f.RoboticsFactory
	} else if ogameID == Shipyard.ID {
		return f.Shipyard
	} else if ogameID == ResearchLab.ID {
		return f.ResearchLab
	} else if ogameID == AllianceDepot.ID {
		return f.AllianceDepot
	} else if ogameID == MissileSilo.ID {
		return f.MissileSilo
	} else if ogameID == NaniteFactory.ID {
		return f.NaniteFactory
	} else if ogameID == Terraformer.ID {
		return f.Terraformer
	} else if ogameID == SpaceDock.ID {
		return f.SpaceDock
	}
	return 0
}

func (f Facilities) String() string {
	return "\n" +
		"RoboticsFactory: " + strconv.Itoa(f.RoboticsFactory) + "\n" +
		"       Shipyard: " + strconv.Itoa(f.Shipyard) + "\n" +
		"   Research Lab: " + strconv.Itoa(f.ResearchLab) + "\n" +
		" Alliance Depot: " + strconv.Itoa(f.AllianceDepot) + "\n" +
		"   Missile Silo: " + strconv.Itoa(f.MissileSilo) + "\n" +
		" Nanite Factory: " + strconv.Itoa(f.NaniteFactory) + "\n" +
		"    Terraformer: " + strconv.Itoa(f.Terraformer) + "\n" +
		"     Space Dock: " + strconv.Itoa(f.SpaceDock)
}
