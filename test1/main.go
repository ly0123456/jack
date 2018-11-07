package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
)

func main() {
	cli := ClI{}

	input := bufio.NewScanner(os.Stdin)

	Welcome()
	cli.Help()

	for {
		fmt.Printf("\ncmd>")
		if input.Scan() {
			s := input.Text()
			cmds := strings.Fields(s)
			if len(cmds) != 0 {
				cli.Run(cmds)
				fmt.Printf("+++++++++++++++++++++++++++++++\n")
			}

		}

	}
}