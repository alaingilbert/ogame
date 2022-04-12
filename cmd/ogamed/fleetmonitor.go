package main

import (
	"github.com/faunX/ogame"
)

func Monitor(bot *ogame.OGame, fleetUpdate chan bool) {

	var fleetCache []ogame.Fleet

	for {

		select {
		case <-fleetUpdate:
			f, _ := bot.GetFleets()
			AddNewFleet(f, fleetCache)
		}
	}
}

func AddNewFleet(fleetOnline []ogame.Fleet, fleetsCache []ogame.Fleet) bool {
	for _, v := range fleetOnline {
		var exists bool
		var updated bool
		var updateKey int64
		for k, w := range fleetsCache {
			if v.ID == w.ID {
				if v.ReturnFlight != w.ReturnFlight {
					updated = true
					updateKey = int64(k)
				}
				exists = true
			}

		}
		if !exists {
			fleetsCache = append(fleetsCache, v)
		}
		if updated {
			fleetsCache[updateKey] = v
		}
	}
	return true
}
