package yandexgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Conn struct {
	API_key  string
	API_url  string
	FolderID string
}

func NewConnection(key, url, folderID string) *Conn {
	return &Conn{API_key: key, API_url: url, FolderID: folderID}
}

type Options struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
}

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type GPTRequest struct {
	ModelUri          string    `json:"modelUri"`
	CompletionOptions Options   `json:"completionOptions"`
	Messages          []Message `json:"messages"`
}

func (s Conn) Promt(request, data string) (string, error) {

	uri := "gpt://" + s.FolderID + "/yandexgpt-lite/latest"

	gptRequest := GPTRequest{
		ModelUri: uri,
		CompletionOptions: Options{
			Stream:      false,
			Temperature: 1.0,
			MaxTokens:   2000,
		},
		Messages: []Message{
			{
				Role: "system",
				Text: request,
			},
			{
				Role: "user",
				Text: data,
			},
		},
	}

	jsonData, err := json.MarshalIndent(gptRequest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling error: %w", err)
	}

	payload := []byte(string(jsonData))

	req, err := http.NewRequest("POST", s.API_url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", s.API_key))
	req.Header.Set("x-folder-id", s.FolderID)
	req.Header.Set("x-data-logging-enabled", "false")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing POST request: %w", err)
	}
	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)

	var gptResponse map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&gptResponse); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	var result string

	respMap := gptResponse["result"].(map[string]interface{})
	alternatives := respMap["alternatives"].([]interface{})
	for _, v := range alternatives {
		alter := v.(map[string]interface{})
		message := alter["message"].(map[string]interface{})
		text := message["text"].(string)

		result = text
	}

	return result, nil
}
