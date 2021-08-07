package ogame

import (
	"encoding/xml"
	"net/http"
	"strconv"
)

type XMLPlayer struct {
	PlayerID int64  `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	Status   string `xml:"status,attr"`
	Alliance string `xml:"alliance,attr"`
}

// Players represent api result https://s157-ru.ogame.gameforge.com/api/players.xml
type XMLPlayers struct {
	Timestamp int64       `xml:"timestamp,attr"`
	ServerID  string      `xml:"serverId,attr"`
	Player    []XMLPlayer `xml:"player"`
}

// gets the server data from xml api
func (b *OGame) getPlayers() (XMLPlayers, error) {
	var players XMLPlayers
	req, err := http.NewRequest("GET", "https://s"+strconv.FormatInt(b.server.Number, 10)+"-"+b.server.Language+".ogame.gameforge.com/api/players.xml", nil)
	if err != nil {
		return players, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req = req.WithContext(b.ctx)
	resp, err := b.Client.Do(req)
	if err != nil {
		return players, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := wrapperReadBody(b, resp)
	if err != nil {
		return players, err
	}
	b.bytesUploaded += req.ContentLength
	if err := xml.Unmarshal(by, &players); err != nil {
		return players, err
	}
	return players, nil
}

func (b *OGame) GetPlayers() (XMLPlayers, error) {
	return b.getPlayers()
}

type Universe struct {
	Timestamp int64  `xml:"timestamp,attr"`
	ServerID  string `xml:"serverId,attr"`
	Planet    []struct {
		PlanetID int64  `xml:"id,attr"`
		PlayerID int64  `xml:"player,attr"`
		Name     string `xml:"name,attr"`
		Coord    string `xml:"coords,attr"`
		Moon     struct {
			MoonID int64  `xml:"id,attr"`
			Name   string `xml:"name,attr"`
			Size   int64  `xml:"size,attr"`
		} `xml:"moon"`
	} `xml:"planet"`
}

// gets the universe data from xml api
func (b *OGame) getUnivers() (Universe, error) {
	var universe Universe
	req, err := http.NewRequest("GET", "https://s"+strconv.FormatInt(b.server.Number, 10)+"-"+b.server.Language+".ogame.gameforge.com/api/universe.xml", nil)
	if err != nil {
		return universe, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req = req.WithContext(b.ctx)
	resp, err := b.Client.Do(req)
	if err != nil {
		return universe, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := wrapperReadBody(b, resp)
	if err != nil {
		return universe, err
	}
	b.bytesUploaded += req.ContentLength
	if err := xml.Unmarshal(by, &universe); err != nil {
		return universe, err
	}
	return universe, nil
}

func (b *OGame) GetUniverse() (Universe, error) {
	return b.getUnivers()
}

type PlayerData struct {
	Username  string
	ID        string
	Positions struct {
		Position []struct {
			Type     int64 `xml:"type,attr"`
			Position int64 `xml:",chardata"`
			Score    int64 `xml:"score,attr"`
			Ships    int64 `xml:"ships,attr"`
		} `xml:"position"`
	} `xml:"positions"`

	Planets struct {
		Planet []struct {
			PlanetID int64  `xml:"id,attr"`
			Name     string `xml:"name,attr"`
			Coords   string `xml:"coords,attr"`
		} `xml:"planet"`
	} `xml:"planets"`
	Alliance struct {
		ID   string `xml:"id,attr"`
		Name string `xml:"name"`
		Tag  string `xml:"tag"`
	} `xml:"alliance"`
}

// gets the PlayerData data from xml api
func (b *OGame) getPlayerData(playerID int64) (PlayerData, error) {
	var playerData PlayerData
	req, err := http.NewRequest("GET", "https://s"+strconv.FormatInt(b.server.Number, 10)+"-"+b.server.Language+".ogame.gameforge.com/api/playerData.xml?id="+strconv.FormatInt(playerID, 10), nil)
	if err != nil {
		return playerData, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req = req.WithContext(b.ctx)
	resp, err := b.Client.Do(req)
	if err != nil {
		return playerData, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := wrapperReadBody(b, resp)
	if err != nil {
		return playerData, err
	}
	b.bytesUploaded += req.ContentLength
	if err := xml.Unmarshal(by, &playerData); err != nil {
		return playerData, err
	}
	return playerData, nil
}

func (b *OGame) GetPlayerData(playerID int64) (PlayerData, error) {
	return b.getPlayerData(playerID)
}

// TOTP
func (b *OGame) RegisterTOTP() {

}
