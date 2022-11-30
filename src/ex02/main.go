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
	arguments := stdinRead[1:]
	command := exec.Command(stdinRead[0], arguments...)
	stdout, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(stdout))
}
