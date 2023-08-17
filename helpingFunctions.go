package memoryProfiling

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	SpeedDial = 100
	Workers   = 8
)

func Loader(output chan string) {
	folderPath := "inputData"

	files, err := ReadFilesFromFolder(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Files in the folder:")
	for _, file := range files {
		output <- file
		fmt.Println(file)
	}
	close(output)
}

func Worker(files chan string, output chan string) {
	for fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

		i := 0
		for scanner.Scan() {
			if i%SpeedDial == 0 {
				allCombinations := GeneratePermutations([]rune(scanner.Text()))
				for index := range allCombinations {
					go SimulateHighMemoryUsage(allCombinations[index])
					output <- allCombinations[index]
					i++
				}
			} else {
				output <- scanner.Text()
			}
		}
	}
}

func ReadFilesFromFolder(folderPath string) ([]string, error) {
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

func WriteToFile(filename string, content string) {
	file, err := os.Create(filename)
	IsFatal(err)
	defer file.Close()
	_, err = file.WriteString(content)
	IsFatal(err)
}

func IsFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GeneratePermutations(chars []rune) []string {
	if len(chars) == 0 {
		return []string{""}
	}

	perms := []string{}
	for i, char := range chars {
		remainingChars := make([]rune, len(chars)-1)
		copy(remainingChars, chars[:i])
		copy(remainingChars[i:], chars[i+1:])

		// Recursively generate permutations for remaining characters
		subPerms := GeneratePermutations(remainingChars)

		// Append the selected character to each sub-permutation
		for _, subPerm := range subPerms {
			perms = append(perms, string(char)+subPerm)
		}
	}

	return perms
}

func SimulateHighMemoryUsage(input string) {
	memoryHog := ""

	for i := 0; i < 750; i++ {
		memoryHog += input
	}
}
