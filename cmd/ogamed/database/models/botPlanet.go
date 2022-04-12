package models

import "github.com/faunX/ogame"

type BotPlanet struct {
	ID                       uint
	BotID                    uint
	Bot                      Bot `gorm:"foreignKey:BotID"`
	ogame.Planet             `gorm:"embedded"`
	ogame.ResourcesBuildings `gorm:"embedded"`
	ogame.Facilities         `gorm:"embedded"`
	ogame.ShipsInfos         `gorm:"embedded"`
	ogame.DefensesInfos      `gorm:"embedded"`
	ogame.Resources          `gorm:"embedded"`
	ogame.ResourceSettings   `gorm:"embedded"`
	//ogame.EspionageReport      `gorm:"embedded"`
	Idx                        uint
	BrainMode                  uint
	EvacuationMode             uint
	Producer                   uint
	Exporter                   uint
	Importer                   uint
	RepatriateActive           bool
	RepatriateMinimumDeuterium uint
	RepatriateMetal            uint
	RepatriateCrystal          uint
	RepatriateDeuterium        uint
}
