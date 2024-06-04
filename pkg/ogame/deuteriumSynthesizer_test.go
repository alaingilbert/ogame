package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeuteriumSynthesizerPrice(t *testing.T) {
	ds := newDeuteriumSynthesizer()

	price := ds.GetPrice(1, LfBonuses{})
	assert.Equal(t, Resources{Metal: 225, Crystal: 75}, price)

	price = ds.GetPrice(2, LfBonuses{})
	assert.Equal(t, Resources{Metal: 337, Crystal: 112}, price)

	price = ds.GetPrice(3, LfBonuses{})
	assert.Equal(t, Resources{Metal: 506, Crystal: 168}, price)

	price = ds.GetPrice(4, LfBonuses{})
	assert.Equal(t, Resources{Metal: 759, Crystal: 253}, price)

	price = ds.GetPrice(5, LfBonuses{})
	assert.Equal(t, Resources{Metal: 1139, Crystal: 379}, price)

	price = ds.GetPrice(11, LfBonuses{})
	assert.Equal(t, Resources{Metal: 12974, Crystal: 4324}, price)
}

func TestDeuteriumSynthesizerConstructionTime(t *testing.T) {
	ds := newDeuteriumSynthesizer()
	assert.Equal(t, 1845*time.Second, ds.ConstructionTime(9, 6, Facilities{}, LfBonuses{}, NoClass, false))
}

func TestDeuteriumSynthesizer_Production(t *testing.T) {
	ds := newDeuteriumSynthesizer()
	assert.Equal(t, int64(40699), ds.Production(7, (-23+17)/2, 1, 1, 15, 28))
}

func TestDeuteriumSynthesizer_EnergyConsumption(t *testing.T) {
	ds := newDeuteriumSynthesizer()
	assert.Equal(t, int64(6198), ds.EnergyConsumption(26))
}
