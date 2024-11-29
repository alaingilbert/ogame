package v12_0_0

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	"math"
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
