package main

import "fmt"

func main() {
	fmt.Println("Welcome to Catan!\n")

	printBoard(*GenerateBoard())
}
