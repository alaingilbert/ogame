package parser

import (
	"github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/extractor/v7"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func MustReadFile(p string) []byte {
	pageHTMLBytes, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return pageHTMLBytes
}

func parserErr(_ any, err error) error {
	return err
}

func TestParsePageContent(t *testing.T) {
	assert.NoError(t, parserErr(ParsePage[DefensesPage](v6.NewExtractor(), MustReadFile("../../samples/unversioned/defence.html"))))
	assert.NoError(t, parserErr(ParsePage[DefensesPage](v7.NewExtractor(), MustReadFile("../../samples/v7/defenses.html"))))
	assert.Error(t, parserErr(ParsePage[DefensesPage](v7.NewExtractor(), MustReadFile("../../samples/v7/overview.html"))))
}

func TestParsePageGetDoc(t *testing.T) {
	p, _ := ParsePage[DefensesPage](v7.NewExtractor(), MustReadFile("../../samples/v7/defenses.html"))
	assert.Nil(t, p.doc)
	p.GetDoc()
	assert.NotNil(t, p.doc)
}
