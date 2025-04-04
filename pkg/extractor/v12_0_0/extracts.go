package v12_0_0

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"math"
	"regexp"
	"strings"
	"time"
)

func extractServerTimeFromDoc(doc *goquery.Document, clock clockwork.Clock) (time.Time, error) {
	txt := doc.Find("div.OGameClock").First().Text()
	serverTime, err := time.Parse("02.01.2006 15:04:05", txt)
	if err != nil {
		return time.Time{}, err
	}

	u1 := clock.Now().UTC().Unix()
	u2 := serverTime.Unix()
	n := int(math.Round(float64(u2-u1)/900)) * 900 // u2-u1 should be close to 0, round to nearest 15min difference

	serverTime = serverTime.Add(time.Duration(-n) * time.Second).In(time.FixedZone("OGT", n))

	return serverTime, nil
}

func extractHighscoreFromDoc(doc *goquery.Document) (out ogame.Highscore, err error) {
	s := doc.Selection
	isFullPage := doc.Find("#stat_list_content").Size() == 1
	if isFullPage {
		s = doc.Find("#stat_list_content")
	}

	script := s.Find("script").First().Text()
	m := regexp.MustCompile(`var site = (\d+);`).FindStringSubmatch(script)
	if len(m) != 2 {
		return out, errors.New("failed to find site")
	}
	out.CurrPage = utils.DoParseI64(m[1])

	m = regexp.MustCompile(`var currentCategory = (\d+);`).FindStringSubmatch(script)
	if len(m) != 2 {
		return out, errors.New("failed to find currentCategory")
	}
	out.Category = utils.DoParseI64(m[1])

	m = regexp.MustCompile(`var currentType = (\d+);`).FindStringSubmatch(script)
	if len(m) != 2 {
		return out, errors.New("failed to find currentType")
	}
	out.Type = utils.DoParseI64(m[1])

	changeSiteSize := s.Find("select.changeSite option").Size()
	out.NbPage = utils.MaxInt(int64(changeSiteSize)-1, 0)

	for _, s := range s.Find("#ranks tbody tr").EachIter() {
		p := ogame.HighscorePlayer{}
		p.Position = utils.DoParseI64(strings.TrimSpace(s.Find("td.position").Text()))
		p.ID = utils.DoParseI64(strings.TrimPrefix(s.AttrOr("id", "position0"), "position"))
		p.Name = strings.TrimSpace(s.Find("span.playername").Text())
		tdName := s.Find("td.name")
		allyTag := tdName.Find("span.ally-tag")
		if allyTag != nil {
			href := allyTag.Find("a").AttrOr("href", "")
			m := regexp.MustCompile(`allianceId=(\d+)`).FindStringSubmatch(href)
			if len(m) == 2 {
				p.AllianceID = utils.DoParseI64(m[1])
			}
			allyTag.Remove()
		}
		href := tdName.Find("a").AttrOr("href", "")
		m := regexp.MustCompile(`galaxy=(\d+)&system=(\d+)&position=(\d+)`).FindStringSubmatch(href)
		if len(m) != 4 {
			continue
		}
		p.Homeworld.Type = ogame.PlanetType
		p.Homeworld.Galaxy = utils.DoParseI64(m[1])
		p.Homeworld.System = utils.DoParseI64(m[2])
		p.Homeworld.Position = utils.DoParseI64(m[3])
		honorScoreSpan := s.Find("span.honorScore span")
		if honorScoreSpan == nil {
			continue
		}
		p.HonourPoints = utils.ParseInt(strings.TrimSpace(honorScoreSpan.Text()))
		p.Score = utils.ParseInt(strings.TrimSpace(s.Find("td.score").Text()))
		shipsRgx := regexp.MustCompile(`([\d\.]+)`)
		shipsTitle := strings.TrimSpace(s.Find("td.score").AttrOr("title", "0"))
		shipsParts := shipsRgx.FindStringSubmatch(shipsTitle)
		if len(shipsParts) == 2 {
			p.Ships = utils.ParseInt(shipsParts[1])
		}
		out.Players = append(out.Players, p)
	}

	return
}
