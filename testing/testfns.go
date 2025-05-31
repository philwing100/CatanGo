// testfns.go
package testing

import (
	"bufio"
	"fmt"
	"io"
)

func promptInt(prompt string, r io.Reader) int {
	fmt.Print(prompt)
	var num int
	fmt.Fscanln(r, &num)
	return num
}

func waitForEnter(playerName string, r io.Reader) {
	fmt.Printf("%s, press ENTER to roll the die...\n", playerName)
	bufio.NewReader(r).ReadBytes('\n')
}
