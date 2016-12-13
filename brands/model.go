package brands

type berthaBrand struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	ImageURL        string `json:"imageurl"`
	Biography       string `json:"biography"`
	TwitterHandle   string `json:"twitterhandle"`
	FacebookProfile string `json:"facebookprofile"`
	LinkedinProfile string `json:"linkedinprofile"`
	TmeIdentifier   string `json:"tmeidentifier"`
}

type brand struct {
	UUID                   string                 `json:"uuid"`
	PrefLabel              string                 `json:"prefLabel"`
	Type                   string                 `json:"type"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers,omitempty"`
	Aliases                []string               `json:"aliases,omitempty"`
	BirthYear              int                    `json:"birthYear,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	Salutation             string                 `json:"salutation,omitempty"`
	EmailAddress           string                 `json:"emailAddress,omitempty"`
	TwitterHandle          string                 `json:"twitterHandle,omitempty"`
	FacebookProfile        string                 `json:"facebookProfile,omitempty"`
	LinkedinProfile        string                 `json:"linkedinProfile,omitempty"`
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
