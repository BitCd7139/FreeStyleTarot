package prmopt

type CardPrmopt struct {
	Name        string `xml:"name"`
	Meaning     string `xml:"meaning"`
	State       string `xml:"state"`
	Description string `xml:"description"`
	KeyWord     string `xml:"keyword"`
	Arcana      string `xml:"arcana"`
	Element     string `xml:"element"`
	Numerology  string `xml:"numerology"`
	Astrology   string `xml:"astrology"`
}
