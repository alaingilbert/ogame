package ogame

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func MustReadFile(p string) []byte {
	pageHTMLBytes, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return pageHTMLBytes
}

func parserErr(_ any, err error) error {
	return err
}

func TestParseDefensesPageContent(t *testing.T) {
	assert.NoError(t, parserErr(ParseDefensesPageContent(&OGame{extractor: NewExtractorV6()}, MustReadFile("samples/unversioned/defence.html"))))
	assert.NoError(t, parserErr(ParseDefensesPageContent(&OGame{extractor: NewExtractorV7()}, MustReadFile("samples/v7/defenses.html"))))
	assert.Error(t, parserErr(ParseDefensesPageContent(&OGame{extractor: NewExtractorV7()}, MustReadFile("samples/v7/overview.html"))))
}
