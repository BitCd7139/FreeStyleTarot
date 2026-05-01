package prmopt

import (
	"encoding/xml"
)

type FinalPrmopt struct {
	XMLname     xml.Name     `xml:"finalPrmopt"`
	PrePrmopt   string       `xml:"prePrmopt"`
	Question    string       `xml:"question"`
	CardContext []CardPrmopt `xml:"cardContext>cardContext"`
}
