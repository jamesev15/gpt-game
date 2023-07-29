package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

func generateRandomNumber(max int) int {
	return rand.Intn(max)
}

func generateRandomNumberGPT(max int) int {
	numberStr, err := gpt_completion(fmt.Sprintf("Generate a random number between 0 and %v. Only return the number without explanation", max))

	if err != nil {
		log.Fatal(err)
	}

	n, err := strconv.Atoi(numberStr)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error converting number: %v to int", n), err)
		return 50
	}

	return n
}

func randomNumberHandler(isActiveGPT bool) func(int) int {
	if isActiveGPT {
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
	fmt.Println("-- Guessing game using GPT3.5 --")
	fmt.Println("Guess the secret number between 0 and 50.")

	openaiToken := readEnvVar("OPENAI_TOKEN")

	var isActiveGPT bool

	if openaiToken != "" {
		isActiveGPT = true
	}

	handler := randomNumberHandler(isActiveGPT)
	secretNumber := handler(50)

	if isActiveGPT {
		clue, err := gpt_completion(fmt.Sprintf("Tell me a funny clue that helps me to guess the secret number. Just returns the clue without the secret number nor the clue's explanation.: %v", secretNumber))

		if err == nil {
			fmt.Println("Clue: ", clue)
		}

	}

	var guess string

	for {
		fmt.Print("Input your guess: ")
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
