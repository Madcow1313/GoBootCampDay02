package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type myFlags struct {
	slFlag, dFlag, fFlag bool
	ext                  string
}

func printSymLink(f string) {
	link, err := os.Readlink(f)
	if err != nil {
		fmt.Println(f, "->", "[broken]")
	} else {
		_, err = os.Stat(link)
		if err != nil {
			fmt.Println(f, "->", "[broken]")
		} else {
			fmt.Println(f, "->", link)
		}
	}

}

func setIfNoFlags(myFlags *myFlags) {
	if !myFlags.dFlag && !myFlags.fFlag && !myFlags.slFlag {
		myFlags.dFlag = true
		myFlags.fFlag = true
		myFlags.slFlag = true
		myFlags.ext = ""
	}
}

func setFlags(myFlags *myFlags) {
	flag.BoolVar(&myFlags.slFlag, "sl", false, "show symbolic links")
	flag.BoolVar(&myFlags.dFlag, "d", false, "show direcories")
	flag.BoolVar(&myFlags.fFlag, "f", false, "show files")
	flag.StringVar(&myFlags.ext, "ext", "", "show files with certain extension")
}

func main() {
	myFlags := new(myFlags)
	setFlags(myFlags)
	flag.Parse()
	setIfNoFlags(myFlags)
	pathToDir := flag.Args()
	if len(pathToDir) == 0 {
		fmt.Println("Error! No path given")
		os.Exit(-1)
	} else if len(pathToDir) > 1 {
		fmt.Println("Error! Too many arguments")
		os.Exit(-1)
	}
	err := filepath.Walk(pathToDir[0], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			err = nil
		}
		if path == pathToDir[0] {
			return nil
		}
		var fi fs.FileInfo
		fi, err = os.Lstat(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if fi.IsDir() && myFlags.dFlag {
			fmt.Println(path)
		} else if fi.Mode()&os.ModeSymlink != 0 && myFlags.slFlag {
			printSymLink(path)
		} else if myFlags.fFlag {
			if len(myFlags.ext) != 0 {
				fileExtension := filepath.Ext(info.Name())
				if fileExtension == "."+myFlags.ext {
					fmt.Println(path)
				}
			} else if !fi.IsDir() {
				fmt.Println(path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("error is", err)
		err = nil
	}
}
