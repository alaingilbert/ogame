<img src="./logo.png" width="300" />

# OGame automation toolkit

- As a library
- As a service (ogamed)

---

# ogame library

### Verify attack example

```go
package main

import "fmt"
import "os"
import "github.com/alaingilbert/ogame"

func main() {
	universe := os.Getenv("UNIVERSE")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	language := os.Getenv("LANGUAGE")
	bot, _ := ogame.New(universe, username, password, language)
	attacked := bot.IsUnderAttack()
	fmt.Println(attacked) // False
}
```

### Available methods

```go
GetServer() Server
SetUserAgent(newUserAgent string)
ServerURL() string
GetLanguage() string
GetPageContent(url.Values) string
PostPageContent(url.Values, url.Values) string
Login() error
Logout()
GetUsername() string
GetUniverseName() string
GetUniverseSpeed() int
GetUniverseSpeedFleet() int
IsDonutGalaxy() bool
IsDonutSystem() bool
ServerVersion() string
ServerTime() time.Time
IsUnderAttack() bool
GetUserInfos() UserInfos
SendMessage(playerID int, message string) error
GetFleets() []Fleet
CancelFleet(FleetID) error
GetAttacks() []AttackEvent
GalaxyInfos(galaxy, system int) ([]PlanetInfos, error)
GetResearch() Researches
GetCachedPlanets() []Planet
GetPlanets() []Planet
GetPlanetByCoord(Coordinate) (Planet, error)
GetPlanet(PlanetID) (Planet, error)
GetEspionageReportMessages() ([]EspionageReportSummary, error)
GetEspionageReport(msgID int) (EspionageReport, error)
DeleteMessage(msgID int) error
FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int)

// Planet specific functions
GetResourceSettings(PlanetID) (ResourceSettings, error)
SetResourceSettings(PlanetID, ResourceSettings) error
GetResourcesBuildings(PlanetID) (ResourcesBuildings, error)
GetDefense(PlanetID) (DefensesInfos, error)
GetShips(PlanetID) (ShipsInfos, error)
GetFacilities(PlanetID) (Facilities, error)
Build(planetID PlanetID, id ID, nbr int) error
BuildCancelable(PlanetID, ID) error
BuildProduction(planetID PlanetID, id ID, nbr int) error
BuildBuilding(planetID PlanetID, buildingID ID) error
BuildTechnology(planetID PlanetID, technologyID ID) error
BuildDefense(planetID PlanetID, defenseID ID, nbr int) error
BuildShips(planetID PlanetID, shipID ID, nbr int) error
GetProduction(PlanetID) ([]Quantifiable, error)
ConstructionsBeingBuilt(PlanetID) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int)
CancelBuilding(PlanetID) error
CancelResearch(PlanetID) error
GetResources(PlanetID) (Resources, error)
SendFleet(planetID PlanetID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources) (FleetID, error)
GetResourcesProductions(PlanetID) (Resources, error)
```

### Full documentation

[https://godoc.org/github.com/alaingilbert/ogame](https://godoc.org/github.com/alaingilbert/ogame)

---

# ogamed service

Download [ogamed binary here](https://github.com/alaingilbert/ogame/releases)

```
./ogamed --universe=Zibal --username=email@email.com --password=secret --language=en
```

```
$ curl 127.0.0.1:8080/bot/is-under-attack
{"Status":"ok","Code":200,"Message":"","Result":false}

$ curl 127.0.0.1:8080/bot/send-message -d 'playerID=123&message="Sup boi!"'
{"Status":"ok","Code":200,"Message":"","Result":null}

$ curl 127.0.0.1:8080/bot/user-infos
{"Status":"ok","Code":200,"Message":"","Result":{"PlayerID":106734,"PlayerName":"Commodore Nomad","Points":43825,"Rank":1130,"Total":1675,"HonourPoints":0}}
```

```
POST /bot/set-user-agent
POST /bot/server-url
POST /bot/page-content
GET  /bot/login
GET  /bot/logout
GET  /bot/server/speed
GET  /bot/server/version
GET  /bot/server/time
GET  /bot/is-under-attack
GET  /bot/user-infos
POST /bot/send-message
GET  /bot/fleets
POST /bot/fleets/:fleetID/cancel
GET  /bot/attacks
GET  /bot/galaxy-infos/:galaxy/:system
GET  /bot/get-research
GET  /bot/planets
GET  /bot/planets/:galaxy/:system/:position
GET  /bot/planets/:planetID
GET  /bot/planets/:planetID/resource-settings
POST /bot/planets/:planetID/resource-settings
GET  /bot/planets/:planetID/resources-buildings
GET  /bot/planets/:planetID/defence
GET  /bot/planets/:planetID/ships
GET  /bot/planets/:planetID/facilities
POST /bot/planets/:planetID/build/:ogameID/:nbr
POST /bot/planets/:planetID/build/cancelable/:ogameID
POST /bot/planets/:planetID/build/production/:ogameID/:nbr
POST /bot/planets/:planetID/build/building/:ogameID
POST /bot/planets/:planetID/build/technology/:ogameID
POST /bot/planets/:planetID/build/defence/:ogameID/:nbr
POST /bot/planets/:planetID/build/ships/:ogameID/:nbr
GET  /bot/planets/:planetID/production
GET  /bot/planets/:planetID/constructions
POST /bot/planets/:planetID/cancel-building
POST /bot/planets/:planetID/cancel-research
GET  /bot/planets/:planetID/resources
POST /bot/planets/:planetID/send-fleet
```
