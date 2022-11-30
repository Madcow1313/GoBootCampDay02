package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"sync"
)

type myFlags struct {
	lFlag,
	mFlag,
	wFlag bool
}

func setFlags(myFlags *myFlags) {
	flag.BoolVar(&myFlags.lFlag, "l", false, "for counting lines")
	flag.BoolVar(&myFlags.mFlag, "m", false, "for counting characters")
	flag.BoolVar(&myFlags.wFlag, "w", false, "for counting words")
	flag.Parse()
}

func countFlags(myFlags *myFlags) int {
	var counter int
	reflVal := reflect.ValueOf(&myFlags).Elem()
	for i := 0; i < 3; i++ {
		rf := reflect.Indirect(reflVal).Field(i)
		if rf.Bool() {
			counter++
		}
	}
	return counter
}

func countWords(path string, wg *sync.WaitGroup) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		wg.Done()
		log.Println("Error! Couldn't read from file ", path)
		return
	}
	strFile := string(file)
	strSlice := strings.Fields(strFile)
	fmt.Printf("%d\t%s\n", len(strSlice), path)
	wg.Done()
}

func countLines(path string, wg *sync.WaitGroup) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		wg.Done()
		log.Println("Error! Couldn't read from file ", path)
		return
	}
	strFile := string(file)
	strSlice := strings.Split(strFile, "\n")
	fmt.Printf("%d\t%s\n", len(strSlice), path)
	wg.Done()
}

func countCharacters(path string, wg *sync.WaitGroup) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		wg.Done()
		log.Println("Error! Couldn't read from file ", path)
		return
	}
	fmt.Printf("%d\t%s\n", len(file), path)
	wg.Done()
}

func main() {
	myFlags := new(myFlags)
	setFlags(myFlags)
	if countFlags(myFlags) > 1 {
		log.Fatal("Error! Too many flags, one expected")
	}
	files := flag.Args()
	if len(files) < 1 {
		log.Fatal("Error! No file given")
	}
	var pointerToFunction func(string, *sync.WaitGroup)
	if myFlags.lFlag {
		pointerToFunction = countLines
	} else if myFlags.mFlag {
		pointerToFunction = countCharacters
	} else {
		pointerToFunction = countWords
	}

	var wg sync.WaitGroup
	for _, path := range files {
		wg.Add(1)
		go pointerToFunction(path, &wg)
	}
	wg.Wait()
}
