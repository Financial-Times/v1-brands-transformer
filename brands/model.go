package brands

type berthaBrand struct {
	Active         bool   `json:"active"`
	PrefLabel      string `json:"prefLabel"`
	Strapline      string `json:"strapline"`
	ImageURL       string `json:"imageurl"`
	DescriptionXML string `json:"descriptionxml"`
	UUID           string `json:"uuid"`
	ParentUUID     string `json:"parentuuid"`
	TmeIdentifier  string `json:"tmeidentifier"`
}

type brand struct {
	UUID                   string                 `json:"uuid"`
	ParentUUID             string                 `json:"parentUUID,omitempty"`
	PrefLabel              string                 `json:"prefLabel,omitempty"`
	Type                   string                 `json:"type,omitempty"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers,omitempty"`
	Aliases                []string               `json:"aliases,omitempty"`
	Strapline              string                 `json:"strapline,omitempty"`
	Description            string                 `json:"description,omitempty"`
	DescriptionXML         string                 `json:"descriptionXML,omitempty"`
	ImageURL               string                 `json:"_imageUrl,omitempty"`
}

type alternativeIdentifiers struct {
	TME   []string `json:"TME,omitempty"`
	UUIDs []string `json:"uuids,omitempty"`
}

type brandLink struct {
	APIURL string `json:"apiUrl"`
}

type brandUUID struct {
	UUID string `json:"ID"`
}
