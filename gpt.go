package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptPayload struct {
	Model       string       `json:"model"`
	Temperature float64      `json:"temperature"`
	Messages    []GptMessage `json:"messages"`
}

type GptChoice struct {
	Message GptMessage `json:"message"`
}

type GptResponse struct {
	Choices []GptChoice `json:"choices"`
}

const OPENAI_URL = "https://api.openai.com/v1/chat/completions"

func gpt_completion(prompt string) (string, error) {
	gptPayload := GptPayload{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
		Messages: []GptMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(gptPayload)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", OPENAI_URL, bytes.NewBuffer(jsonData))

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", readEnvVar("OPENAI_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("error gpt request")

	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var gptResponse GptResponse
	err = json.Unmarshal(responseBytes, &gptResponse)

	if err != nil {
		return "", err
	}

	return gptResponse.Choices[0].Message.Content, nil
}
