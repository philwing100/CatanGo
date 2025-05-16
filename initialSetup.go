package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func gameStart() {

}

func cliGameInitialize() int {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Welcome to Catan!")
		fmt.Print("Please enter the number of players (3 or 4): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}

		input = strings.TrimSpace(input)
		playerNum, err := strconv.Atoi(input)
		if err != nil || (playerNum != 3 && playerNum != 4) {
			fmt.Println("Invalid input. Please enter either 3 or 4.")
			continue
		}

		return playerNum
	}
}
