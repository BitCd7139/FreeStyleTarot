package request

type Predict struct {
	Question       string           `json:"question"`
	IntentSummary  string           `json:"intent_summary,omitempty"`
	Clarifications []Clarification  `json:"clarifications,omitempty"`
	CardSize       CardSize         `json:"cardSize"`
	Model          string           `json:"model"`
	Cards          []CardInfo       `json:"cards"`
	FreestyleMode  bool             `json:"freestylemode,omitempty"`
	CustomAPI      *CustomAPIConfig `json:"custom_api,omitempty"`
}

// CustomAPIConfig holds user-provided LLM API configuration.
// When set, the backend uses this instead of the global LLM.
type CustomAPIConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
	Format  string `json:"format"` // "openai" or "anthropic"
	// Per-stage overrides. Zero values fall back to server defaults.
	StageParams *StageParams `json:"stage_params,omitempty"`
}

// StageParams allows the user to override max_tokens and temperature
// for each agent stage when using a custom API.
type StageParams struct {
	IntentClarifier *StageParam `json:"intent_clarifier,omitempty"`
	SpreadAnalyst    *StageParam `json:"spread_analyst,omitempty"`
	PersonaAdvisor   *StageParam `json:"persona_advisor,omitempty"`
}

// StageParam holds tunable generation parameters for one stage.
type StageParam struct {
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type CardSize struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

type CardInfo struct {
	Order       int     `json:"order"`
	Name        string  `json:"name"`
	X           float32 `json:"x"`
	Y           float32 `json:"y"`
	Orientation string  `json:"orientation"`
	Meaning     string  `json:"meaning"`
}
