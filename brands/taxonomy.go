package brands

type taxonomy struct {
	Terms []term `xml:"term"`
}

//TODO revise fields for brands - Also need labels to come through too
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
