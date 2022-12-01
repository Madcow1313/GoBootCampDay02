package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

func packFiles(file string, wg *sync.WaitGroup) (string, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	mtime := fi.ModTime().Unix()
	str := strconv.Itoa(int(mtime))
	tarName := fi.Name() + "_" + str + ".tar.gz"
	outFile, err := os.Create(tarName)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	gw := gzip.NewWriter(outFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	err = addToArchive(tw, file)
	if err != nil {
		return "", err
	}
	return tarName, nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = filename
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}
	return nil
}

func startGoRoutine(filename string, wg *sync.WaitGroup, aFlag string) {
	str, err := packFiles(filename, wg)
	if err != nil {
		log.Println(err, filename)
	} else if len(aFlag) != 0 {
		err = os.Rename(str, aFlag+"/"+str)
		if err != nil {
			log.Println(err, aFlag)
		}
	}
	wg.Done()
}

func main() {
	var aFlag string
	flag.StringVar(&aFlag, "a", "", "where to put archives")
	flag.Parse()
	files := flag.Args()
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go startGoRoutine(file, &wg, aFlag)
	}
	wg.Wait()
}
