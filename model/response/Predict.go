package response

type Predict struct {
	Answer string `json:"answer"`
	Code   int    `json:"code"`
}
