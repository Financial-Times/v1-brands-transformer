package brands

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	taxonomyName = "taxonomy_name"
)

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
