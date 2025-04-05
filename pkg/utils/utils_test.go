package utils

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

func TestParseInt(t *testing.T) {
	assert.Equal(t, int64(1234567890), ParseInt("1.234.567.890"))
	assert.Equal(t, int64(0), ParseInt(""))
}

func TestParseInt2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../samples/unversioned/deathstar_price.html")
	doc := First(goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes)))
	title := doc.Find("li.metal").AttrOr("title", "")
	metalStr := regexp.MustCompile(`([\d.]+)`).FindStringSubmatch(title)[1]
	metal := ParseInt(metalStr)
	assert.Equal(t, int64(5000000), metal)

	pageHTMLBytes, _ = os.ReadFile("../../samples/unversioned/mrd_price.html")
	doc = First(goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes)))
	title = doc.Find("li.metal").AttrOr("title", "")
	metalStr = regexp.MustCompile(`([\d.]+)`).FindStringSubmatch(title)[1]
	metal = ParseInt(metalStr)
	assert.Equal(t, int64(1555733200), metal)
}

func TestToInt(t *testing.T) {
	assert.Equal(t, 1234567890, ToInt([]byte("1234567890")))
}

func TestMinInt(t *testing.T) {
	assert.Equal(t, int64(2), MinInt(5, 2, 3))
}

func TestI64Ptr(t *testing.T) {
	v := int64(6)
	assert.Equal(t, &v, I64Ptr(6))
}

func TestDefault(t *testing.T) {
	assert.Equal(t, true, Default((*bool)(nil), true))
	assert.Equal(t, false, Default((*bool)(nil), false))
	assert.Equal(t, true, Default(Ptr(true), false))
	assert.Equal(t, false, Default(Ptr(false), true))
	assert.Equal(t, 1, Default((*int)(nil), 1))
	assert.Equal(t, 2, Default(Ptr(2), 1))
}
