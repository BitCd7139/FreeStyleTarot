package service

const (
	APIKey  = "你的DeepSeek_API_KEY"
	BaseURL = "https://api.deepseek.com/chat/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

//func callDeepSeek(fullPrompt string) (string, error) {
//	reqBody := ChatRequest{
//		Model: "deepseek-chat", // 或者 deepseek-reasoning
//		Messages: []Message{
//			{"role": "system", "content": "You are a helpful assistant."}, // 实际操作中可将固定Prompt放这
//			{"role": "user", "content": fullPrompt},
//		},
//		Stream: false,
//	}
//
//	jsonData, _ := json.Marshal(reqBody)
//
//	req, _ := http.NewRequest("POST", BaseURL, bytes.NewBuffer(jsonData))
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+APIKey)
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	body, _ := ioutil.ReadAll(resp.Body)
//
//	return string(body), nil
//}
