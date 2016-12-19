package brands

import (
	"encoding/base64"
	"encoding/xml"

	"github.com/pborman/uuid"
)

// BrandTransformer struct
type BrandTransformer struct {
}

// UnMarshallTaxonomy - unmarshal the XML of a taxonomy
func (*BrandTransformer) UnMarshallTaxonomy(contents []byte) ([]interface{}, error) {
	t := taxonomy{}
	err := xml.Unmarshal(contents, &t)
	if err != nil {
		return nil, err
	}
	interfaces := make([]interface{}, len(t.Terms))
	for i, d := range t.Terms {
		interfaces[i] = d
	}
	return interfaces, nil
}

// UnMarshallTerm - unmarshal the XML of a TME term
func (*BrandTransformer) UnMarshallTerm(content []byte) (interface{}, error) {
	dummyTerm := term{}
	err := xml.Unmarshal(content, &dummyTerm)
	if err != nil {
		return term{}, err
	}
	return dummyTerm, nil
}

func transformBrand(tmeTerm term, taxonomyName string) brand {
	tmeIdentifier := buildTmeIdentifier(tmeTerm.RawID, taxonomyName)
	brandUUID := uuid.NewMD5(uuid.UUID{}, []byte(tmeIdentifier)).String()
	uuidList := []string{brandUUID}
	if val, ok := berthaUUIDmap()[tmeIdentifier]; ok {
		brandUUID = val
		uuidList = append(uuidList, val)
	}

	aliasList := buildAliasList(tmeTerm.Aliases, tmeTerm.CanonicalName)
	return brand{
		UUID:      brandUUID,
		PrefLabel: tmeTerm.CanonicalName,
		AlternativeIdentifiers: alternativeIdentifiers{
			TME:   []string{tmeIdentifier},
			UUIDs: removeDuplicates(uuidList),
		},
		Type:    "Brand",
		Aliases: aliasList,
	}
}

func buildTmeIdentifier(rawID string, tmeTermTaxonomyName string) string {
	id := base64.StdEncoding.EncodeToString([]byte(rawID))
	taxonomyName := base64.StdEncoding.EncodeToString([]byte(tmeTermTaxonomyName))
	return id + "-" + taxonomyName
}

func removeDuplicates(slice []string) []string {
	newSlice := []string{}
	seen := make(map[string]bool)
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			newSlice = append(newSlice, v)
			seen[v] = true
		}
	}
	return newSlice
}

func buildAliasList(aList aliases, canonicalName string) []string {
	aliasList := make([]string, len(aList.Alias))
	for k, v := range aList.Alias {
		aliasList[k] = v.Name
	}
	aliasList = append(aliasList, canonicalName)
	aliasList = removeDuplicates(aliasList)
	return aliasList
}
