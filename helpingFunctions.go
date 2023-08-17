package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func loader(output chan string) {
	folderPath := "inputData"

	// Read the files in the folder
	files, err := readFilesFromFolder(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	// Print the list of files
	fmt.Println("Files in the folder:")
	for _, file := range files {
		output <- file
		fmt.Println(file)
	}
	close(output)
}

func worker(files chan string, output chan string) {
	for fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			output <- scanner.Text()
		}
		time.Sleep(time.Millisecond * 1)
	}
}

func readFilesFromFolder(folderPath string) ([]string, error) {
	var files []string

	entries, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// Check if the entry is a file (not a directory)
		if entry.Mode().IsRegular() {
			files = append(files, filepath.Join(folderPath, entry.Name()))
		}
	}

	return files, nil
}

func writeToFile(filename string, content string) {
	file, err := os.Create(filename)
	isFatal(err)
	defer file.Close()
	_, err = file.WriteString(content)
	isFatal(err)
}

func isFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
