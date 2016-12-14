package brands

type taxonomy struct {
	Terms []term `xml:"term"`
}

type term struct {
	CanonicalName string  `xml:"name"`
	RawID         string  `xml:"id"`
	Aliases       aliases `xml:"variations"`
}

type aliases struct {
	Alias []alias `xml:"variation"`
}

type alias struct {
	Name string `xml:"name"`
}
