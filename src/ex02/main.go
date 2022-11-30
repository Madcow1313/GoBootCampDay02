package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var stdinRead []string
	for scanner.Scan() {
		text := scanner.Text()
		splitted := strings.Fields(text)
		stdinRead = append(stdinRead, splitted...)
	}
	arguments := os.Args[2:]
	arguments = append(arguments, stdinRead...)
	command := exec.Command(os.Args[1], arguments...)
	stdout2, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(stdout2))
}
