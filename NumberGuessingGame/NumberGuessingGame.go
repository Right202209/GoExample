package main

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomNumber(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max) + 1
}
func bj(num int, ans int) string {
	if num < ans {
		return "greater"
	} else {
		return "less"
	}
}
func main() {
	fmt.Println("Welcome to the NumberGuessingGame!")
	fmt.Println("I'm thinking of a number between 1 and 100.\nYou have some chances to guess the correct number.")
	fmt.Println("Please select the difficulty level:\n1. Easy (10 chances)\n2. Medium (5 chances)\n3. Hard (3 chances)")

	var levelN int
	var level string
	var chance int
	chance = 1

	fmt.Print("Enter your choice:")
	fmt.Scan(&levelN)
	switch levelN {
	case 1:
		level = "Easy"
		chance = 10
	case 2:
		level = "Medium"
		chance = 5
	case 3:
		level = "Hard"
		chance = 3
	default:
		fmt.Println("Please select 1 - 3")
	}
	fmt.Printf("Great! You have selected the %s difficulty level.", level)
	fmt.Println("Let's start the game!")
	Num := RandomNumber(100)
	fmt.Printf("answer:%d", Num)
	var guess int

	for i := 0; i < chance; {
		fmt.Scan(&guess)
		if guess == Num {
			fmt.Printf("Congratulations! You guessed the correct number in %d attempts.", i+1)
			return
		} else {
			var B string
			B = bj(guess, Num)
			fmt.Printf("Incorrect! The number is %s than %d.\n", B, guess)
			i++
		}
	}
	fmt.Println("Game Over")
	return
}
