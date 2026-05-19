package prompt

import (
	"encoding/xml"
)

type FinalPrompt struct {
	XMLname     xml.Name     `xml:"finalPrompt"`
	PrePrompt   string       `xml:"prePrompt"`
	Question    string       `xml:"question"`
	CardContext []CardPrompt `xml:"cardContext>cardContext"`
}
