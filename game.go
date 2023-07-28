package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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

func readEnvVar(env string) string {
	return os.Getenv(env)
}

func generateRandomNumber(max int) int {
	return rand.Intn(max)
}

func generateRandomNumberGPT(max int) int {
	gptPayload := GptPayload{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
		Messages: []GptMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf("Generate a random number between 0 and %v. Only return the number without explanation", max),
			},
		},
	}

	jsonData, err := json.Marshal(gptPayload)

	if err != nil {
		fmt.Println("Internal Error converting Payload for GPT", err)
		return 50
	}

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", readEnvVar("OPENAI_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return 50
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error in response")
		return 50
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
		return 50
	}

	var gptResponse GptResponse
	err = json.Unmarshal(responseBytes, &gptResponse)

	if err != nil {
		fmt.Println("Error parsing GPT response")
		return 50
	}

	numberStr := gptResponse.Choices[0].Message.Content

	n, err := strconv.Atoi(numberStr)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error converting number: %v to int", n), err)
		return 50
	}

	return n
}

func randomNumberHandler() func(int) int {
	openaiToken := readEnvVar("OPENAI_TOKEN")

	if openaiToken != "" {
		fmt.Println("OPENAI TOKEN found!")
		fmt.Println("Generating random number with GPT 3.5")
		return generateRandomNumberGPT
	} else {
		fmt.Println("OPENAI TOKEN not found")
		fmt.Println("Generating random number with cpu")
		return generateRandomNumber
	}
}

func main() {
	fmt.Println("Guessing game using GPT3.5")

	handler := randomNumberHandler()
	secretNumber := handler(100)

	var guess string

	for {
		fmt.Println("Input your guess: ")
		_, err := fmt.Scanln(&guess)

		if err != nil {
			fmt.Println("Error reading from CLI")
			log.Fatal(err)
		}

		guessNumber, err := strconv.Atoi(guess)

		if err != nil {
			fmt.Println("It must be a number. Try again")
			continue
		}

		if guessNumber > secretNumber {
			fmt.Println("To big")
		} else if guessNumber < secretNumber {
			fmt.Println("To less")
		} else {
			fmt.Println("Win!")
			break
		}
	}

}
