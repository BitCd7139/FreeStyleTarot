package response

import "FreeStyleTarot/model/request"

type Clarify struct {
	NeedsClarification bool                      `json:"needs_clarification"`
	IntentSummary      string                    `json:"intent_summary,omitempty"`
	Questions          []request.ClarifyQuestion `json:"questions"`
}
