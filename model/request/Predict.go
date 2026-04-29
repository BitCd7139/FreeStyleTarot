package request

type Predict struct {
	Question string     `json:"question"`
	CardSize CardSize   `json:"cardSize"`
	Cards    []CardInfo `json:"cards"`
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
