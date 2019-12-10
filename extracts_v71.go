package ogame

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type resourcesRespV71 struct {
	Resources struct {
		Metal struct {
			Amount  float64 `json:"amount"`
			Storage float64 `json:"storage"`
			Tooltip string  `json:"tooltip"`
		} `json:"metal"`
		Crystal struct {
			Amount  float64 `json:"amount"`
			Storage float64 `json:"storage"`
			Tooltip string  `json:"tooltip"`
		} `json:"crystal"`
		Deuterium struct {
			Amount  float64 `json:"amount"`
			Storage float64 `json:"storage"`
			Tooltip string  `json:"tooltip"`
		} `json:"deuterium"`
		Energy struct {
			Amount  float64 `json:"amount"`
			Tooltip string  `json:"tooltip"`
		} `json:"energy"`
		Darkmatter struct {
			Amount  float64 `json:"amount"`
			Tooltip string  `json:"tooltip"`
		} `json:"darkmatter"`
	} `json:"resources"`
	HonorScore int64 `json:"honorScore"`
	Techs      struct {
		Num1 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"1"`
		Num2 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"2"`
		Num3 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"3"`
		Num4 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"4"`
		Num12 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"12"`
		Num212 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"212"`
		Num217 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"217"`
	} `json:"techs"`
}

func extractResourcesDetailsV71(pageHTML []byte) (out ResourcesDetails, err error) {
	var res resourcesRespV71
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		fmt.Println("CALSS", err)
		if isLogged(pageHTML) {
			return out, ErrInvalidPlanetID
		}
		return
	}
	out.Metal.Available = int64(res.Resources.Metal.Amount)
	out.Metal.StorageCapacity = int64(res.Resources.Metal.Storage)
	out.Crystal.Available = int64(res.Resources.Crystal.Amount)
	out.Crystal.StorageCapacity = int64(res.Resources.Crystal.Storage)
	out.Deuterium.Available = int64(res.Resources.Deuterium.Amount)
	out.Deuterium.StorageCapacity = int64(res.Resources.Deuterium.Storage)
	out.Energy.Available = int64(res.Resources.Energy.Amount)
	out.Darkmatter.Available = int64(res.Resources.Darkmatter.Amount)
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Metal.Tooltip))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Crystal.Tooltip))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Deuterium.Tooltip))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Darkmatter.Tooltip))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Energy.Tooltip))
	out.Metal.CurrentProduction = ParseInt(metalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Crystal.CurrentProduction = ParseInt(crystalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Deuterium.CurrentProduction = ParseInt(deuteriumDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Energy.CurrentProduction = ParseInt(energyDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Energy.Consumption = ParseInt(energyDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Darkmatter.Purchased = ParseInt(darkmatterDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Darkmatter.Found = ParseInt(darkmatterDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	return
}
