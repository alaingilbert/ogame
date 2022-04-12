package main

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/faunX/ogame"
)

type Localization struct {
	XMLName                   xml.Name `xml:"localization"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	Timestamp                 string   `xml:"timestamp,attr"`
	ServerId                  string   `xml:"serverId,attr"`
	Techs                     struct {
		Text string `xml:",chardata"`
		Name []struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"name"`
	} `xml:"techs"`
	Missions struct {
		Text string `xml:",chardata"`
		Name []struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"name"`
	} `xml:"missions"`
}

var localization Localization

func LocalizeObject(name string) ogame.BaseOgameObj {
	var obj ogame.BaseOgameObj
	for _, t := range localization.Techs.Name {
		if strings.ToLower(t.Text) == strings.ToLower(name) {
			id, _ := strconv.ParseInt(t.ID, 10, 64)
			obj = ogame.Objs.ByID(ogame.ID(id))
		}
	}
	return obj
}
