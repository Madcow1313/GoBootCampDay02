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
	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewReader(stdout)
	newArgs := os.Args[2:]
	command = exec.Command(os.Args[1], newArgs...)
	oldstdiin := os.Stdin
	command.Stdin = buf
	stdoutstr, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdin = oldstdiin
	fmt.Println(stdoutstr)
}
