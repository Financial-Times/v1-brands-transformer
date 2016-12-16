package brands

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	taxonomyName = "taxonomy_name"
)

func TestDodgyUUIDTransformBrand(t *testing.T) {
	testTerm := term{
		CanonicalName: "Business blog",
		RawID:         "Brands_86",
		Aliases: aliases{
			Alias: []alias{
				{Name: "Business Blog"},
				{Name: "The Blog of Business"},
			}},
	}
	brandTaxonomy := "Brands"

	tfp := transformBrand(testTerm, brandTaxonomy)
	assert.EqualValues(t, []string{"0312776d-bac4-3118-bc76-b93b2cd3f1ba", "fd4459b2-cc4e-4ec8-9853-c5238eb860fb"}, tfp.AlternativeIdentifiers.UUIDs)
	assert.Equal(t, "fd4459b2-cc4e-4ec8-9853-c5238eb860fb", tfp.UUID)
}

func TestTransformBrand(t *testing.T) {
	testTerm := term{
		CanonicalName: "Bob",
		RawID:         "bob",
		Aliases: aliases{
			Alias: []alias{
				{Name: "B"},
				{Name: "b"},
			}},
	}
	tfp := transformBrand(testTerm, taxonomyName)
	log.Infof("got brand %v", tfp)
	assert.NotNil(t, tfp)
	assert.EqualValues(t, []string{"B", "b", "Bob"}, tfp.Aliases)
	assert.Equal(t, "0e86d39b-8320-3a98-a87a-ff35d2cb04b9", tfp.UUID)
	assert.Equal(t, "Bob", tfp.PrefLabel)
}

func TestDeduplication(t *testing.T) {
	input := []string{"a", "b", "b", "c", "d", "d"}
	expectedOutput := []string{"a", "b", "c", "d"}

	actualOutput := removeDuplicates(input)
	assert.Equal(t, expectedOutput, actualOutput)
}

func TestBuildAliasList(t *testing.T) {
	inputAliases := aliases{Alias: []alias{
		alias{Name: "A"},
		alias{Name: "B"},
		alias{Name: "C"},
	}}
	inputName := "C"

	expectedOutput := []string{"A", "B", "C"}

	actualOutput := buildAliasList(inputAliases, inputName)
	assert.EqualValues(t, expectedOutput, actualOutput)
}
