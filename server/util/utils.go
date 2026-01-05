package utils

import (
	"bufio"
	"log"
	"os"
)

type Note struct {
	NoteCount int
	Notes     []string
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	Check(err)
	if os.IsNotExist(err) {
		return nil
	}
	Check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	Check(scanner.Err())
	return lines
}
