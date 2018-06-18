package astrophysics

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestCost(t *testing.T) {
	a := New()
	assert.Equal(t, ogame.Resources{Metal: 7000, Crystal: 14000, Deuterium: 7000}, a.GetPrice(2))
}
