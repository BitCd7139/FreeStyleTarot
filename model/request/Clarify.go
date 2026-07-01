package request

type Clarify struct {
	Question  string           `json:"question"`
	Cards     []CardInfo       `json:"cards"`
	CustomAPI *CustomAPIConfig `json:"custom_api,omitempty"`
}

type ClarifyOption struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type ClarifyQuestion struct {
	ID                string          `json:"id"`
	Text              string          `json:"text"`
	Options           []ClarifyOption `json:"options"`
	AllowCustom       bool            `json:"allow_custom"`
	MultiSelect       bool            `json:"multi_select,omitempty"`
	CustomPlaceholder string          `json:"custom_placeholder,omitempty"`
}

type Clarification struct {
	QuestionID string   `json:"questionId"`
	Question   string   `json:"question"`
	OptionID   string   `json:"optionId,omitempty"`
	OptionIDs  []string `json:"optionIds,omitempty"`
	Answer     string   `json:"answer"`
}
