package wrapper

import (
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// GetFleetSpeedForMission ...
func GetFleetSpeedForMission(serverData gameforge.ServerData, missionID ogame.MissionID) int64 {
	if missionID == ogame.Attack ||
		missionID == ogame.GroupedAttack ||
		missionID == ogame.Destroy ||
		missionID == ogame.MissileAttack ||
		missionID == ogame.RecycleDebrisField {
		return serverData.SpeedFleetWar
	}
	return serverData.SpeedFleetPeaceful
}
