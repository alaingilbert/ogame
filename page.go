package ogame

import "time"

type Page struct {
	b       *OGame
	content []byte
}

type EventListAjaxPage struct{ Page }
type MissileAttackLayerAjaxPage struct{ Page }
type FetchTechsAjaxPage struct{ Page }
type RocketlayerAjaxPage struct{ Page }

type FullPage struct{ Page }
type OverviewPage struct{ FullPage }
type SuppliesPage struct{ FullPage }
type ResourcesSettingsPage struct{ FullPage }

//type FacilitiesPageContent struct{ FullPageContent }
//type TraderOverviewPageContent struct{ FullPageContent }
//type TraderResourcesPageContent struct{ FullPageContent }

type ResearchPage struct{ FullPage }
type FacilitiesPage struct{ FullPage }
type ShipyardPage struct{ FullPage }
type DefensesPage struct{ FullPage }

//type FleetDispatchPageContent struct{ FullPageContent }

type MovementPage struct{ FullPage }

//type GalaxyPageContent struct{ FullPageContent }
//type AlliancePageContent struct{ FullPageContent }
//type PremiumPageContent struct{ FullPageContent }
//type ShopPageContent struct{ FullPageContent }
//type MessagesPageContent struct{ FullPageContent }
//type ChatPageContent struct{ FullPageContent }
//type CharacterClassSelectionPageContent struct{ FullPageContent }
//type BuddiesPageContent struct{ FullPageContent }
//type HighScorePageContent struct{ FullPageContent }

func (p RocketlayerAjaxPage) ExtractDestroyRockets() (int64, int64, string, error) {
	return p.b.extractor.ExtractDestroyRockets(p.content)
}

func (p FullPage) ExtractServerTime() (time.Time, error) {
	return p.b.extractor.ExtractServerTime(p.content)
}

func (p FullPage) ExtractPlanets() []Planet {
	return p.b.extractor.ExtractPlanets(p.content, p.b)
}

func (p FullPage) ExtractPlanet(v any) (Planet, error) {
	return p.b.extractor.ExtractPlanet(p.content, p.b, v)
}

func (p FullPage) ExtractMoons() []Moon {
	return p.b.extractor.ExtractMoons(p.content, p.b)
}

func (p FullPage) ExtractMoon(v any) (Moon, error) {
	return p.b.extractor.ExtractMoon(p.content, p.b, v)
}

func (p FullPage) ExtractCelestials() ([]Celestial, error) {
	return p.b.extractor.ExtractCelestials(p.content, p.b)
}

func (p FullPage) ExtractCelestial(v any) (Celestial, error) {
	return p.b.extractor.ExtractCelestial(p.content, p.b, v)
}

func (p ResearchPage) ExtractResearch() Researches {
	return p.b.extractor.ExtractResearch(p.content)
}

func (p SuppliesPage) ExtractResourcesBuildings() (ResourcesBuildings, error) {
	return p.b.extractor.ExtractResourcesBuildings(p.content)
}

func (p DefensesPage) ExtractDefense() (DefensesInfos, error) {
	return p.b.extractor.ExtractDefense(p.content)
}

func (p OverviewPage) ExtractDMCosts() (DMCosts, error) {
	return p.b.extractor.ExtractDMCosts(p.content)
}

func (p OverviewPage) ExtractConstructions() (ID, int64, ID, int64) {
	return p.b.extractor.ExtractConstructions(p.content)
}

func (p OverviewPage) ExtractUserInfos() (UserInfos, error) {
	return p.b.extractor.ExtractUserInfos(p.content, p.b.language)
}

func (p OverviewPage) ExtractCancelResearchInfos() (token string, techID, listID int64, err error) {
	return p.b.extractor.ExtractCancelResearchInfos(p.content)
}

func (p OverviewPage) ExtractCancelBuildingInfos() (token string, techID, listID int64, err error) {
	return p.b.extractor.ExtractCancelBuildingInfos(p.content)
}

func (p FacilitiesPage) ExtractFacilities() (Facilities, error) {
	return p.b.extractor.ExtractFacilities(p.content)
}

func (p ShipyardPage) ExtractProduction() ([]Quantifiable, int64, error) {
	return p.b.extractor.ExtractProduction(p.content)
}

func (p ShipyardPage) ExtractShips() (ShipsInfos, error) {
	return p.b.extractor.ExtractShips(p.content)
}

func (p ResourcesSettingsPage) ExtractResourceSettings() (ResourceSettings, error) {
	return p.b.extractor.ExtractResourceSettings(p.content)
}

func (p MovementPage) ExtractFleets() []Fleet {
	return p.b.extractor.ExtractFleets(p.content, p.b.location)
}

func (p MovementPage) ExtractSlots() Slots {
	return p.b.extractor.ExtractSlots(p.content)
}

func (p MovementPage) ExtractCancelFleetToken(fleetID FleetID) (string, error) {
	return p.b.extractor.ExtractCancelFleetToken(p.content, fleetID)
}

func (p EventListAjaxPage) ExtractAttacks(ownCoords []Coordinate) ([]AttackEvent, error) {
	return p.b.extractor.ExtractAttacks(p.content, ownCoords)
}

func (p MissileAttackLayerAjaxPage) ExtractIPM() (int64, int64, string) {
	return p.b.extractor.ExtractIPM(p.content)
}

func (p FetchTechsAjaxPage) ExtractTechs() (ResourcesBuildings, Facilities, ShipsInfos, DefensesInfos, Researches, error) {
	return p.b.extractor.ExtractTechs(p.content)
}

type FullPagePages interface {
	OverviewPage |
		SuppliesPage |
		ResourcesSettingsPage |
		FacilitiesPage |
		//TraderOverviewPageContent |
		//TraderResourcesPageContent |
		ResearchPage |
		ShipyardPage |
		DefensesPage |
		//FleetDispatchPageContent |
		MovementPage
	//GalaxyPageContent |
	//AlliancePageContent |
	//PremiumPageContent |
	//ShopPageContent |
	//MessagesPageContent |
	//ChatPageContent |
	//CharacterClassSelectionPageContent |
	//BuddiesPageContent |
	//HighScorePageContent
}

type AjaxPagePages interface {
	EventListAjaxPage |
		MissileAttackLayerAjaxPage |
		FetchTechsAjaxPage |
		RocketlayerAjaxPage
}
