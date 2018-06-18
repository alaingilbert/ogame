package rocketLauncher

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// RocketLauncher ...
type RocketLauncher struct {
	baseDefense.BaseDefense
}

// New ...
func New() *RocketLauncher {
	d := new(RocketLauncher)
	d.OGameID = 401
	d.Price = ogame.Resources{Metal: 2000}
	d.StructuralIntegrity = 2000
	d.ShieldPower = 20
	d.WeaponPower = 80
	d.RapidfireFrom = map[ogame.ID]int{ogame.Bomber: 20, ogame.Cruiser: 10, ogame.Deathstar: 200}
	d.Requirements = map[ogame.ID]int{ogame.Shipyard: 1}
	return d
}
