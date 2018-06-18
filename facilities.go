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
	if ogameID == RoboticsFactory {
		return f.RoboticsFactory
	} else if ogameID == Shipyard {
		return f.Shipyard
	} else if ogameID == ResearchLab {
		return f.ResearchLab
	} else if ogameID == AllianceDepot {
		return f.AllianceDepot
	} else if ogameID == MissileSilo {
		return f.MissileSilo
	} else if ogameID == NaniteFactory {
		return f.NaniteFactory
	} else if ogameID == Terraformer {
		return f.Terraformer
	} else if ogameID == SpaceDock {
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
